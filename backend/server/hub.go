package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Hub struct {
	bus      *BusServer
	clients  map[string]*Client
	clientsM sync.RWMutex
}

func NewHub(bus *BusServer) *Hub {
	return &Hub{bus: bus, clients: make(map[string]*Client)}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		ID:            genUUID(),
		Conn:          conn,
		Send:          make(chan *Envelope, 128),
		Subscriptions: make(map[string]bool),
	}

	h.clientsM.Lock()
	h.clients[client.ID] = client
	h.clientsM.Unlock()

	// register client with bus
	h.bus.RegisterClient(client)

	// start read/write pumps
	go h.readPump(client)
	go h.writePump(client)
}

func (h *Hub) readPump(c *Client) {
	defer func() {
		// cleanup
		h.clientsM.Lock()
		delete(h.clients, c.ID)
		h.clientsM.Unlock()

		h.bus.UnregisterClient(c.ID)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(8192)
	// No initial read deadline, rely on ping/pong for keepalive

	for {
		var env Envelope
		if err := c.Conn.ReadJSON(&env); err != nil {
			log.Println("read json err:", err)
			break
		}

		// Let the bus handle all the routing logic
		env.From = c.ID
		h.bus.Publish(&env)
	}
}

func (h *Hub) writePump(c *Client) {
	for env := range c.Send {
		c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.Conn.WriteJSON(env); err != nil {
			log.Println("write json err:", err)
			return
		}
	}
}
