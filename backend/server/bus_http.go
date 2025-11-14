package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewServerMux(hub *Hub) *http.ServeMux {
	muxRoot := http.NewServeMux()
	router := mux.NewRouter()
	router.HandleFunc("/ws", hub.ServeWS)

	// Serve static files
	fs := http.FileServer(http.Dir("../frontend/dist"))
	router.PathPrefix("/").Handler(fs)

	muxRoot.Handle("/", router)

	// health
	muxRoot.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	log.Println("HTTP routes registered")
	return muxRoot
}
