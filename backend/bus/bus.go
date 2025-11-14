package bus

import (
	"sync"
)

// Message represents a message on the bus
type Message struct {
	Topic   string
	Payload interface{}
	// Add other fields from your spec like ID, From, To, Token, etc.
}

// Client represents a client connected to the bus
type Client struct {
	ID      string
	Receive chan Message
}

// Topic represents a topic on the bus
type Topic struct {
	Name        string
	subscribers map[string]*Client
	mu          sync.RWMutex
}

// Bus is the central message bus
type Bus struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

// NewBus creates a new Bus
func NewBus() *Bus {
	return &Bus{
		topics: make(map[string]*Topic),
	}
}

// Subscribe allows a client to subscribe to a topic
func (b *Bus) Subscribe(topicName string, client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic, exists := b.topics[topicName]
	if !exists {
		topic = &Topic{
			Name:        topicName,
			subscribers: make(map[string]*Client),
		}
		b.topics[topicName] = topic
	}

	topic.mu.Lock()
	defer topic.mu.Unlock()
	topic.subscribers[client.ID] = client
}

// Unsubscribe allows a client to unsubscribe from a topic
func (b *Bus) Unsubscribe(topicName string, client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic, exists := b.topics[topicName]
	if !exists {
		return
	}

	topic.mu.Lock()
	defer topic.mu.Unlock()
	delete(topic.subscribers, client.ID)
}

// Publish sends a message to all subscribers of a topic
func (b *Bus) Publish(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topic, exists := b.topics[msg.Topic]
	if !exists {
		return
	}

	topic.mu.RLock()
	defer topic.mu.RUnlock()

	for _, client := range topic.subscribers {
		go func(c *Client) {
			c.Receive <- msg
		}(client)
	}
}
