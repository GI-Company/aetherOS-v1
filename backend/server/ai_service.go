package server

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

// AIService handles interactions with the generative AI model.
type AIService struct {
	bus         *BusServer
	genaiClient *genai.Client
}

// NewAIService creates and initializes a new AIService.
func NewAIService(bus *BusServer) *AIService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: GEMINI_API_KEY environment variable not set. AI service will be disabled.")
		return nil // Or return a disabled service
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
    APIKey: apiKey,
  })
	if err != nil {
		log.Printf("Failed to create GenAI client: %v", err)
		return nil
	}

	ai := &AIService{
		bus:         bus,
		genaiClient: client,
	}

	// Subscribe to AI requests from the bus
	bus.SubscribeServer("ai:request", ai.handleAIRequest)

	log.Println("AI service initialized successfully.")
	return ai
}

// handleAIRequest processes AI requests received from the bus.
func (ai *AIService) handleAIRequest(env *Envelope) {
	// For now, we'll just log the request.
	// In the future, this will call the Gemini API.
	log.Printf("AI request received: %v", env.Payload)
}
