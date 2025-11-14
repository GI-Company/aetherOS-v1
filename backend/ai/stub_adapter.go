package ai

import (
	"context"
	"time"
)

type StubAdapter struct{}

func NewStubAdapter() *StubAdapter { return &StubAdapter{} }

func (s *StubAdapter) GenerateText(ctx context.Context, prompt string) (AIResult, error) {
	// simple deterministic stub; replace with real client later
	time.Sleep(50 * time.Millisecond)
	return AIResult{Text: "stub-response: " + prompt}, nil
}
