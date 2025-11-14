package main

import (
	"log"
	"net/http"
	"time"

	"aether/backend/ai"
	"aether/backend/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from the root directory
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	bus := server.NewBusServer()
	cache := server.NewPersistentCache(2048, 30*time.Minute, bus)
	sess := server.NewSessionManager(cache)
	_ = server.NewAppManager(bus, sess)
	_ = server.NewVFSService(bus)
	_ = server.NewWasmRunner(bus)
	_ = server.NewAuthService()
	// bus.Auth = auth

	_, err = ai.NewAIService(bus)
	if err != nil {
		log.Fatalf("Failed to create AI service: %v", err)
	}

	hub := server.NewHub(bus)
	router := server.NewServerMux(hub)

	// Registering the routes from handlers.go
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/apps", appsHandler)
	router.HandleFunc("/apps/{app_id}/instances", createInstanceHandler).Methods("POST")
	router.HandleFunc("/apps/{app_id}/instances/{instance_id}/start", startInstanceHandler).Methods("POST")
	router.HandleFunc("/kv/{key}", kvGetHandler).Methods("GET")

	// Registering the WebSocket handler
	router.HandleFunc("/v1/sessions/{sid}/ws", wsHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("../frontend/dist"))
	router.PathPrefix("/").Handler(fs)

	log.Println("Aether kernel starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
