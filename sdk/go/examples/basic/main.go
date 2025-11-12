package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	aethersdk "aether-sdk-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create client
	client, err := aethersdk.NewClient("http://localhost:8080", "test-jwt-token")
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe and get WS client (bidirectional)
	ws, err := client.Subscribe(ctx, "test-topic")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// reader goroutine
	go func() {
		for env := range ws.Receive() {
			fmt.Printf("Received via WS: topic=%s from=%s payload=%s\n", env.Topic, env.From, string(env.Payload))
		}
		fmt.Println("WS receive channel closed")
	}()

	// send an envelope via WS
	env := &aethersdk.Envelope{
		Topic:   "test-topic",
		Type:    "event",
		Payload: jsonRaw(map[string]string{"hello": "bi-directional"}),
	}
	if err := client.SendEnvelope(ctx, ws, env); err != nil {
		log.Printf("send error: %v", err)
	}

	// also test HTTP publish (unchanged)
	_ = client.Publish(ctx, "test-topic", map[string]string{"hello": "via-http"})

	// wait a bit
	time.Sleep(2 * time.Second)
}

// helper to produce json.RawMessage easily
func jsonRaw(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}
