
// ===========================
// backend/server/hub.go
// ===========================
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default.
		// In a production environment, you should implement a proper origin check.
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	isClosed bool
	mu       sync.Mutex
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	bus        *BusServer
}

func NewHub(bus *BusServer) *Hub {
	hub := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		bus:        bus,
	}
	// Subscribe the Hub to all topics on the bus, so it can forward them.
	// We use a wildcard subscription.
	hub.bus.Subscribe("*", hub.handleBusMessage)
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.mu.Lock()
				if !client.isClosed {
					close(client.send)
					client.isClosed = true
				}
				client.mu.Unlock()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// handleBusMessage is a callback for the BusServer. It forwards messages from the bus
// to the WebSocket clients.
func (h *Hub) handleBusMessage(msg *Message) {
	// Don't forward messages that came from a WebSocket client back to it.
	// This prevents infinite loops. The 'Source' will be the client's connection pointer.
	if msg.Source != nil {
		if _, ok := msg.Source.(*Client); ok {
			return // This message originated from a client, don't broadcast it back.
		}
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshaling bus message to JSON: %v", err)
		return
	}
	h.broadcast <- jsonMsg
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message from client: %v", err)
			continue
		}
		// Set the source of the message to the client pointer.
		// This allows other parts of the system to know the message origin.
		msg.Source = c
		// Publish the message to the bus.
		c.hub.bus.Publish(&msg)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}

// HandleWebSocket handles websocket requests from the peer.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
