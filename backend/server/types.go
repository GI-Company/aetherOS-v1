package main

import (
	"encoding/json"
	"time"
)

// Envelope is the core message used across Aether's broker.
type Envelope struct {
	ID          string            `json:"id"`            // uuid
	From        string            `json:"from,omitempty"`// e.g. "window:1234"
	To          string            `json:"to,omitempty"`  // optional direct-to
	Topic       string            `json:"topic,omitempty"`
	Type        string            `json:"type,omitempty"`     // "rpc","event","stream"
	ContentType string            `json:"contentType,omitempty"`
	Payload     json.RawMessage   `json:"payload,omitempty"`  // raw JSON or base64 binary wrapper
	Meta        map[string]string `json:"meta,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty"`
}

// Simple helper to marshal envelope to JSON bytes
func (e *Envelope) Bytes() []byte {
	b, _ := json.Marshal(e)
	return b
}
