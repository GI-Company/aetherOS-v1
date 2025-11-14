package ai

import (
	"context"
	"log"
	"os"

	"aether/backend/server"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIService handles interactions with the generative AI model.
type AIService struct {
	bus   *server.BusServer
	model *genai.GenerativeModel
}

// NewAIService creates and initializes a new AIService.
func NewAIService(bus *server.BusServer) (*AIService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-1.5-pro")

	ai := &AIService{
		bus:   bus,
		model: model,
	}

	bus.SubscribeServer("ai:generate:text", ai.handleGenerateText)
	log.Println("AI Service Started")
	return ai, nil
}

func (ai *AIService) handleGenerateText(env *server.Envelope) {
	// 1. Extract prompt and messageId from the payload
	var prompt string
	var messageID interface{}

	if p, ok := env.Payload["prompt"].(string); ok {
		prompt = p
	} else {
		log.Println("handleGenerateText: missing 'prompt' in payload")
		return
	}

	messageID, _ = env.Payload["messageId"] // It's ok if it's missing

	// 2. Generate content using the AI model
	resp, err := ai.model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		log.Printf("handleGenerateText: error generating content: %v", err)
		return
	}

	// 3. Extract the generated text
	var generatedText string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					generatedText += string(txt)
				}
			}
		}
	}

	// 4. Publish the response
	responsePayload := map[string]interface{}{"text": generatedText}
	if messageID != nil {
		responsePayload["messageId"] = messageID
	}

	responseEnv := &server.Envelope{
		Topic:   "ai:generated:text",
		Payload: responsePayload,
	}
	ai.bus.Publish(responseEnv)
}
