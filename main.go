
package main

import (
	"log"
	"net/http"

	"aether/backend/server"
	"github.com/gorilla/mux"
)

func main() {
	bus := server.NewBusServer()
	hub := server.NewHub(bus)
	go hub.Run()

	vfs := server.NewVFSService(bus)
	sess := server.NewSessionManager(bus, vfs)
	_ = server.NewAppManager(bus, vfs, sess)
	_ = server.NewAIService(bus) // Initialize the AI service

	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.HandleWebSocket(w, r)
	})

	r.HandleFunc("/aetherscript", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement the AetherScript service handler.
		w.WriteHeader(http.StatusNotImplemented)
	})

	log.Println("Starting server on :8080")
	log.Println("Make sure you have set your GEMINI_API_KEY environment variable.")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
