package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewServerMux(hub *Hub) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ws", hub.ServeWS)

	// health
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	log.Println("HTTP routes registered")
	return router
}
