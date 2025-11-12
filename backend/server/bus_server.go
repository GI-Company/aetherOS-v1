package main

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type BusServer struct {
	Broker *Broker
}

func (s *BusServer) handlePublish(w http.ResponseWriter, r *http.Request) {
	claims := FromContextClaims(r.Context())
	// example: check sub claim exists
	if claimsMap, ok := claims.(jwt.MapClaims); ok {
		if _, ok := claimsMap["sub"].(string); !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}
	}

	var env Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Broker.Publish(&env)
	w.WriteHeader(http.StatusAccepted)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

func (s *BusServer) handleWSSubscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		conn.Close()
		return
	}

	// Create a new client and subscribe it to the topic
	client := NewClient(conn)
	s.Broker.Subscribe(client, topic)
}
