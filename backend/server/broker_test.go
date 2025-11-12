package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestBroker(t *testing.T) {
	t.Run("Subscribe and Publish", func(t *testing.T) {
		b := NewBroker(context.Background())
		defer b.Shutdown(context.Background())

		sub, err := b.Subscribe("test.topic", 10, nil)
		if err != nil {
			t.Fatalf("Subscribe() error = %v", err)
		}

		env := &Envelope{
			Topic:   "test.topic",
			Payload: []byte(`{"msg":"hello"}`),
		}

		if err := b.Publish(env); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}

		select {
		case receivedEnv := <-sub.Ch:
			if !cmp.Equal(env, receivedEnv) {
				t.Errorf("Received envelope does not match published envelope. Diff: %s", cmp.Diff(env, receivedEnv))
			}
		case <-time.After(1 * time.Second):
			t.Fatal("timed out waiting for message")
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		b := NewBroker(context.Background())
		defer b.Shutdown(context.Background())

		sub, _ := b.Subscribe("test.topic", 10, nil)
		b.Unsubscribe(sub)

		// After unsubscribe, the channel should be closed
		select {
		case _, ok := <-sub.Ch:
			if ok {
				t.Error("subscriber channel should be closed after unsubscribe, but it was not")
			}
		case <-time.After(1 * time.Second):
			t.Fatal("timed out waiting for channel to close")
		}
	})

	t.Run("Direct Message", func(t *testing.T) {
		b := NewBroker(context.Background())
		defer b.Shutdown(context.Background())

		sub, _ := b.Subscribe("test.topic", 10, nil)

		env := &Envelope{
			To:      sub.ID,
			Payload: []byte(`{"msg":"direct"}`),
		}

		if err := b.Publish(env); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}

		select {
		case receivedEnv := <-sub.Ch:
			if !cmp.Equal(env, receivedEnv) {
				t.Errorf("Received envelope does not match published envelope. Diff: %s", cmp.Diff(env, receivedEnv))
			}
		case <-time.After(1 * time.Second):
			t.Fatal("timed out waiting for direct message")
		}
	})

}
