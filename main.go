package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aether-broker/broker/aether"
	"aether-broker/broker/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	aiModule, err := aether.NewAIModule()
	if err != nil {
		log.Fatalf("failed to create AI module: %v", err)
	}
	// The new AI client does not have a Close() method, so we remove the defer call.

	broker := aether.NewBroker()
	go broker.Run()

	r := mux.NewRouter()
	server.RegisterBusRoutes(r, broker)

	// Add the new AI endpoint
	r.HandleFunc("/v1/ai/generate", func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Prompt string `json:"prompt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := aiModule.GenerateText(reqBody.Prompt)
		if err != nil {
			http.Error(w, "failed to generate text", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"text": resp})
	}).Methods("POST")

	// Add the new multimodal AI endpoint
	r.HandleFunc("/v1/ai/multimodal", func(w http.ResponseWriter, r *http.Request) {
		// Implement multimodal logic here
	}).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("starting server on :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
		}

	log.Println("server gracefully stopped")
}
