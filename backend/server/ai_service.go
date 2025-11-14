
// ===============================
// backend/server/ai_service.go
// ===============================
package server

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

// AIService handles interactions with the generative AI model.
type AIService struct {
	bus       *BusServer
	genaiClient *genai.Client
}

// NewAIService creates a new AIService.
func NewAIService(bus *BusServer) *AIService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set.")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	ai := &AIService{
		bus:       bus,
		genaiClient: client,
	}

	bus.Subscribe("ai:generate", ai.handleGenerate)
	return ai
}

// handleGenerate handles requests to generate content.
func (ai *AIService) handleGenerate(msg *Message) {
	prompt, ok := msg.Payload["prompt"].(string)
	if !ok {
		ai.bus.Reply(msg, map[string]interface{}{"error": "prompt not found in payload"})
		return
	}

	model := ai.genaiClient.GenerativeModel("gemini-1.5-pro")
	ctx := context.Background()
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error generating content: %v", err)
		ai.bus.Reply(msg, map[string]interface{}{"error": "failed to generate content"})
		return
	}

	// Assuming printResponse is a function that can format the response.
	// For now, we'll just send the raw response back.
	// In a real application, you would parse the response and send a structured message.
	ai.bus.Reply(msg, map[string]interface{}{"response": resp})
}
