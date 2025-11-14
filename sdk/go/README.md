# Aether SDK for Go

This is the Go SDK for interacting with an Aether broker.

## Installation

```bash
go get aether-sdk-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aethersdk "aether-sdk-go"
)

func main() {
	brokerURL := os.Getenv("AETHER_BROKER_URL")
	if brokerURL == "" {
		brokerURL = "http://localhost:8080"
	}

	// In a real application, you would obtain a JWT from your auth provider.
	// For this example, we'll generate a dummy one.
	testToken, err := aethersdk.NewJWT("test-user", time.Hour)
	if err != nil {
		log.Fatalf("error creating test token: %v", err)
	}

	client, err := aethersdk.NewClient(brokerURL, testToken)
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to a topic
	msgs, err := client.Subscribe(ctx, "test-topic")
	if err != nil {
		log.Fatalf("error subscribing: %v", err)
	}

	// Start a goroutine to listen for messages
	go func() {
		for msg := range msgs {
			fmt.Printf("Received message: %+v\n", msg)
		}
	}()

	// Publish a message
	payload := map[string]interface{}{"hello": "world"}
	if err := client.Publish(ctx, "test-topic", payload); err != nil {
		log.Fatalf("error publishing: %v", err)
	}

	// Wait for a bit to receive the message
	time.Sleep(2 * time.Second)
}

```

## Troubleshooting

If you encounter issues while running the application, here are a few common problems and their solutions:

*   **Dependency Issues:** If you see errors related to missing packages, especially after adding a new dependency, run `go mod vendor` in the `backend` directory. This ensures that the `vendor` folder is in sync with `go.mod`.
*   **`.env` File Not Found:** If the application can't find your `.env` file (and therefore your `GEMINI_API_KEY`), make sure the `godotenv.Load()` function in `backend/main.go` is pointing to the correct location. If your `.env` file is in the project root, the line should be `godotenv.Load("../.env")`.
*   **Port Conflict ("address already in use"):** This error means another process is using the port the application needs (usually 8080). You can find and stop the conflicting process with the following command:
    ```bash
    kill $(lsof -t -i:8080)
    ```
