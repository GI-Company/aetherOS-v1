
// =================================
// backend/server/bus_server.go
// =================================
package server

import (
	"context"
	"sync"
	"time"
)

type BusServer struct {
	mu          sync.RWMutex
	subscribers map[string][]func(*Message)
	requests    map[string]chan *Message
}

func NewBusServer() *BusServer {
	return &BusServer{
		subscribers: make(map[string][]func(*Message)),
		requests:    make(map[string]chan *Message),
	}
}

func (b *BusServer) Subscribe(topic string, handler func(*Message)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], handler)
}

func (b *BusServer) Publish(msg *Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if subscribers, ok := b.subscribers[msg.Topic]; ok {
		for _, handler := range subscribers {
			go handler(msg)
		}
	}
}

func (b *BusServer) PublishSync(topic string, payload map[string]interface{}, timeout time.Duration) *Message {
	reqID := genUUID()
	replyChan := make(chan *Message, 1)
	b.mu.Lock()
	b.requests[reqID] = replyChan
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		delete(b.requests, reqID)
		b.mu.Unlock()
	}()

	msg := &Message{
		Topic:   topic,
		Payload: payload,
		Token:   reqID, // Using token to carry the request ID
	}
	b.Publish(msg)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case reply := <-replyChan:
		return reply
	case <-ctx.Done():
		return nil // Timeout
	}
}

func (b *BusServer) Reply(originalMsg *Message, replyPayload map[string]interface{}) {
	reqID := originalMsg.Token
	b.mu.RLock()
	replyChan, ok := b.requests[reqID]
	b.mu.RUnlock()

	if ok {
		replyChan <- &Message{
			Topic:   "reply:" + originalMsg.Topic,
			Payload: replyPayload,
		}
	}
}
