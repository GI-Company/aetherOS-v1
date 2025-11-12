package aether

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// AIModule handles interactions with the Gemini API.
type AIModule struct {
	client *genai.Client
}

// NewAIModule creates a new AI module.
func NewAIModule() (*AIModule, error) {
	ctx := context.Background()
	// The client automatically uses the GEMINI_API_KEY environment variable.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &AIModule{client: client}, nil
}

// GenerateText generates text from a given prompt.
func (m *AIModule) GenerateText(prompt string) (string, error) {
	ctx := context.Background()
	// Switched to gemini-1.5-flash as a more standard model choice.
	result, err := m.client.Models.GenerateContent(
		ctx,
		"gemini-1.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	text := result.Text()
	if text == "" {
		return "", fmt.Errorf("no text content found in response")
	}

	return text, nil
}

// The new client from genai.NewClient does not have a Close() method.
// Connection management is handled automatically.
