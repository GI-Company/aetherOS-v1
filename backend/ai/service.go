package ai

import (
	"context"
	"log"

	"aether-broker/backend/bus"
	"google.golang.org/genai"
	"google.golang.org/api/option"
)

// AIService handles all AI-related operations
type AIService struct {
	bus    *bus.Bus
	client *bus.Client
	genAI  *genai.GenerativeModel
}

// NewAIService creates a new AIService
func NewAIService(bus *bus.Bus, apiKey string) (*AIService, error) {
	ctx := context.Background()
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := genaiClient.GenerativeModel("gemini-1.5-pro")

	return &AIService{
		bus:    bus,
		client: &bus.Client{ID: "ai-service", Receive: make(chan bus.Message, 128)},
		genAI:  model,
	}, nil
}

// Start begins the AI service's message listening loop
func (s *AIService) Start() {
	log.Println("Starting AI Service")
	s.bus.Subscribe("ai:generate", s.client)
	go s.listen()
}

func (s *AIService) listen() {
	for msg := range s.client.Receive {
		if msg.Topic == "ai:generate" {
			go s.handleGenerate(msg)
		}
	}
}

func (s *AIService) handleGenerate(msg bus.Message) {
	prompt, ok := msg.Payload.(string)
	if !ok {
		log.Println("AI service received invalid payload")
		return
	}

	log.Printf("AI Service received prompt: %s", prompt)

	ctx := context.Background()
	resp, err := s.genAI.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("error generating content: %v", err)
		// Optionally, publish an error message back
		return
	}

	// Assuming the response can be represented as text
	var responseText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					responseText += string(txt)
				}
			}
		}
	}

	s.bus.Publish(bus.Message{Topic: "ai:generate:resp", Payload: responseText})
}
