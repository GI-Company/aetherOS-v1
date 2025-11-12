package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type BusServer struct {
	Broker *Broker
}

// RegisterBusRoutes attaches routes to mux.Router
func RegisterBusRoutes(r *mux.Router, b *Broker) {
	s := &BusServer{Broker: b}
	api := r.PathPrefix("/v1/bus").Subrouter()
	api.Use(JWTMiddleware)
	api.HandleFunc("/publish", s.handlePublish).Methods("POST")
	api.HandleFunc("/subscribe", s.handleWSSubscribe) // upgrade to ws (GET or POST)
}

// handlePublish accepts a JSON Envelope and publishes it.
func (s *BusServer) handlePublish(w http.ResponseWriter, r *http.Request) {
	var env Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "invalid envelope: "+err.Error(), http.StatusBadRequest)
		return
	}
	if env.ID == "" {
		env.ID = uuidSafe()
	}
	if err := s.Broker.Publish(&env); err != nil {
		http.Error(w, "publish error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": env.ID})
}

// WS upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // adjust origin policy
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// handleWSSubscribe upgrades and subscribes to topic(s).
// query parameters:
//  - topic (required): topic name to subscribe
//  - sid (optional): session id for metadata
func (s *BusServer) handleWSSubscribe(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	topic := q.Get("topic")
	if topic == "" {
		http.Error(w, "topic required", http.StatusBadRequest)
		return
	}
	sid := q.Get("sid")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade:", err)
		return
	}
	// subscribe with moderate buffer
	meta := map[string]string{}
	if sid != "" {
		meta["sid"] = sid
	}
	sub, err := s.Broker.Subscribe(topic, 64, meta)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("subscribe error: "+err.Error()))
		_ = conn.Close()
		return
	}
	// consumer goroutine: read messages that client sends (treat incoming JSON as publish)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				// client closed
				return
			}
			switch mt {
			case websocket.TextMessage:
				// treat as Envelope publish JSON
				var env Envelope
				if err := json.Unmarshal(msg, &env); err != nil {
					_ = conn.WriteMessage(websocket.TextMessage, []byte("invalid envelope: "+err.Error()))
					continue
				}
				if env.ID == "" {
					env.ID = uuidSafe()
				}
				if env.Topic == "" {
					// if client didn't supply topic, assume they want to publish to the subscribed topic
					env.Topic = topic
				}
				// set From using sid if present
				if env.From == "" && sid != "" {
					env.From = "session:" + sid
				}
				if err := s.Broker.Publish(&env); err != nil {
					_ = conn.WriteMessage(websocket.TextMessage, []byte("publish error: "+err.Error()))
				}
			case websocket.BinaryMessage:
				// optional: wrap binary frames into an envelope with contentType
				env := Envelope{
					ID:          uuidSafe(),
					Topic:       topic,
					From:        "session:" + sid,
					ContentType: "application/octet-stream",
					Payload:     msg,
					CreatedAt:   time.Now().UTC(),
				}
				_ = s.Broker.Publish(&env)
			}
		}
	}()

	// producer loop: forward broker messages to websocket
	writeTimeout := 5 * time.Second
	for {
		select {
		case env, ok := <-sub.Ch:
			if !ok {
				// broker closed this subscription
				_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "broker closed"))
				_ = conn.Close()
				return
			}
			// send JSON
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteJSON(env); err != nil {
				// client likely disconnected
				_ = conn.Close()
				s.Broker.Unsubscribe(sub)
				return
			}
		case <-ctx.Done():
			// cleanup on reader goroutine done
			s.Broker.Unsubscribe(sub)
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closed"))
			_ = conn.Close()
			return
		case <-s.Broker.ctx.Done():
			// broker shutdown
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown"))
			_ = conn.Close()
			return
		}
	}
}
