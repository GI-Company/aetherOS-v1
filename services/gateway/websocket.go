package gateway

import (
	"encoding/json"
	"log"
	"net/http"

	"aether/services/vfs"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	vfsService *vfs.Service
}

func NewHub(vfsService *vfs.Service) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
		vfsService: vfsService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("error: %v", err)
					delete(h.clients, client)
					client.Close()
				}
			}
		}
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	h.register <- conn
	defer func() { h.unregister <- conn }()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("unmarshal:", err)
			continue
		}

		switch msg["type"] {
		case "vfs:list":
			h.handleVfsList(conn, msg["payload"].(map[string]interface{}))
		case "vfs:create:file":
			h.handleVfsCreateFile(conn, msg["payload"].(map[string]interface{}))
		case "vfs:create:folder":
			h.handleVfsCreateFolder(conn, msg["payload"].(map[string]interface{}))
		case "vfs:delete":
			h.handleVfsDelete(conn, msg["payload"].(map[string]interface{}))
		}
	}
}

func (h *Hub) handleVfsList(conn *websocket.Conn, payload map[string]interface{}) {
	path := payload["path"].(string)
	files, err := h.vfsService.List(path)
	if err != nil {
		log.Println("vfs.List:", err)
		return
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":    "vfs:list:result",
		"payload": map[string]interface{}{"files": files},
	})
	conn.WriteMessage(websocket.TextMessage, response)
}

func (h *Hub) handleVfsCreateFile(conn *websocket.Conn, payload map[string]interface{}) {
	path := payload["path"].(string)
	if err := h.vfsService.CreateFile(path); err != nil {
		log.Println("vfs.CreateFile:", err)
		return
	}
	h.broadcastList(h.vfsService.GetParent(path))
}

func (h *Hub) handleVfsCreateFolder(conn *websocket.Conn, payload map[string]interface{}) {
	path := payload["path"].(string)
	if err := h.vfsService.CreateDirectory(path); err != nil {
		log.Println("vfs.CreateDirectory:", err)
		return
	}
	h.broadcastList(h.vfsService.GetParent(path))
}

func (h *Hub) handleVfsDelete(conn *websocket.Conn, payload map[string]interface{}) {
	path := payload["path"].(string)
	if err := h.vfsService.Delete(path); err != nil {
		log.Println("vfs.Delete:", err)
		return
	}
	h.broadcastList(h.vfsService.GetParent(path))
}

func (h *Hub) broadcastList(path string) {
	files, err := h.vfsService.List(path)
	if err != nil {
		log.Println("vfs.List:", err)
		return
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":    "vfs:list:result",
		"payload": map[string]interface{}{"files": files},
	})
	h.broadcast <- response
}
