package server_test

import (
	"testing"

	"aether/backend/server"
)

func TestBroker(t *testing.T) {
	t.Run("Subscribe and Publish", func(t *testing.T) {
		b := server.NewBusServer()

		env := &server.Envelope{
			Topic:   "test.topic",
			Payload: map[string]interface{}{"msg": "hello"},
		}

		b.Publish(env)
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		// Test unsubscribing
	})

	t.Run("Direct Message", func(t *testing.T) {
		// Test direct messaging
	})
}
