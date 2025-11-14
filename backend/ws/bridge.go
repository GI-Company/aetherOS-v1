package ws

import (
	"log"
	"net/http"

	"aether-broker/backend/bus"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
}

// Bridge handles WebSocket connections and bridges them to the bus
type Bridge struct {
	bus *bus.Bus
}

// NewBridge creates a new Bridge
func NewBridge(bus *bus.Bus) *Bridge {
	return &Bridge{bus: bus}
}

// ServeHTTP handles incoming WebSocket requests
func (b *Bridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade connection: %v", err)
		return
	}

	// Create a new client for this connection
	client := &bus.Client{
		ID:      conn.RemoteAddr().String(), // Use remote address as a unique ID
		Receive: make(chan bus.Message),
	}

	// For now, subscribe the client to all topics. This will be refined later.
	// In a real implementation, the client would send subscription messages.
	b.bus.Subscribe("vfs:write", client)
	b.bus.Subscribe("vfs:read", client)
	b.bus.Subscribe("vfs:list", client)

	go b.readPump(conn, client)
	go b.writePump(conn, client)

	log.Printf("client connected: %s", client.ID)
}

func (b *Bridge) readPump(conn *websocket.Conn, client *bus.Client) {
	defer func() {
		b.bus.Unsubscribe("vfs:write", client) // Clean up subscriptions
		b.bus.Unsubscribe("vfs:read", client)
		b.bus.Unsubscribe("vfs:list", client)
		conn.Close()
		log.Printf("client disconnected: %s", client.ID)
	}()

	for {
		var msg bus.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		b.bus.Publish(msg)
	}
}

func (b *Bridge) writePump(conn *websocket.Conn, client *bus.Client) {
	for msg := range client.Receive {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("error writing json: %v", err)
			break
		}
	}
}
