package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("ws upgrade:", err)
        return
    }
    defer conn.Close()

    vars := mux.Vars(r)
    sid := vars["sid"]
    log.Printf("New WS session %s", sid)

    // Example: read incoming envelopes and dispatch
    for {
        mt, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        if mt == websocket.TextMessage {
            var env map[string]interface{}
            if err := json.Unmarshal(msg, &env); err != nil {
                log.Println("invalid json:", err)
                continue
            }
            go handleClientJSON(context.Background(), sid, env, conn)
        } else {
            // binary frames -> forward to wasm as audio chunk
            go handleClientBinary(context.Background(), sid, msg)
        }
    }
}

func handleClientJSON(ctx context.Context, sid string, env map[string]interface{}, conn *websocket.Conn) {
	log.Printf("Received JSON from %s: %v", sid, env)
	// In a real implementation, this would route to the appropriate instance
	conn.WriteJSON(map[string]string{"status": "received"})
}

func handleClientBinary(ctx context.Context, sid string, msg []byte) {
	log.Printf("Received binary data from %s: %d bytes", sid, len(msg))
	// In a real implementation, this would forward to the wasm instance
}
