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
	// Load .env file from the parent directory
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
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
	mux := server.NewServerMux(hub)

	log.Println("Aether kernel starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
