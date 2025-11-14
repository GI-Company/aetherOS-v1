package ai

import "context"

type AIResult struct {
	Text string
}

type AIAdapter interface {
	GenerateText(ctx context.Context, prompt string) (AIResult, error)
}
