package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// create broker and register bus routes
	broker := NewBroker(context.Background())
	RegisterBusRoutes(r, broker)

	// Register metrics routes
	RegisterMetricsRoutes(r)

	// Serve the frontend
	fs := http.FileServer(http.Dir("../frontend/dist"))
	r.PathPrefix("/").Handler(fs)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
