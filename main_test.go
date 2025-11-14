
package main

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"aether/broker/aether"
	"aether/broker/server"
	aethersdk "aether/sdk"

	"github.com/gorilla/mux"
)

func TestBrokerPubSub(t *testing.T) {
	// Start the broker in a goroutine
	broker := aether.NewBroker()
	go broker.Run()

	r := mux.NewRouter()
	server.RegisterBusRoutes(r, broker)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	defer srv.Shutdown(context.Background())

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	brokerURL := "http://localhost:8080"
	testToken, err := aethersdk.NewJWT("test-user", time.Hour)
	if err != nil {
		t.Fatalf("error creating test token: %v", err)
	}

	client, err := aethersdk.NewClient(brokerURL, testToken)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Subscribe to a topic
	msgs, err := client.Subscribe(ctx, "test-topic")
	if err != nil {
		t.Fatalf("error subscribing: %v", err)
	}

	// Publish a message in a separate goroutine
	go func() {
		time.Sleep(1 * time.Second) // Give the subscriber time to connect
		payload := map[string]string{"message": "hello"}
		if err := client.Publish(ctx, "test-topic", payload); err != nil {
			t.Errorf("error publishing: %v", err)
		}
	}()

	// Wait for the message
	select {
	case env := <-msgs:
		if env.Topic != "test-topic" {
			t.Errorf("got topic %q, want %q", env.Topic, "test-topic")
		}
		// Assuming the payload is a raw JSON message, we'll need to unmarshal it
		var payload map[string]string
		if err := env.UnmarshalPayload(&payload); err != nil {
			t.Fatalf("error unmarshaling payload: %v", err)
		}
		if payload["message"] != "hello" {
			t.Errorf("got message %q, want %q", payload["message"], "hello")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for message")
	}
}
