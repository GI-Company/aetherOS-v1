package server

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Subscriber is an internal handle returned when subscribing.
type Subscriber struct {
	ID       string
	Topic    string
	Ch       chan *Envelope // channel for delivering envelopes
	cancel   context.CancelFunc
	metadata map[string]string
}

// Broker is a simple topic-based pub/sub broker.
// It's memory-resident, concurrency-safe, and supports multiple subscribers per topic.
type Broker struct {
	// topics -> subscriberID -> *Subscriber
	topics map[string]map[string]*Subscriber
	// optional ring buffer per topic could be added later
	mu sync.RWMutex

	// shutdown control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewBroker constructs a broker with a cancellable context.
func NewBroker(parent context.Context) *Broker {
	ctx, cancel := context.WithCancel(parent)
	return &Broker{
		topics: make(map[string]map[string]*Subscriber),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Subscribe creates a subscriber channel for a topic. BufferSize recommended 32.
func (b *Broker) Subscribe(topic string, bufferSize int, metadata map[string]string) (*Subscriber, error) {
	if topic == "" {
		return nil, errors.New("topic required")
	}
	if bufferSize <= 0 {
		bufferSize = 32
	}
	ctx, cancel := context.WithCancel(b.ctx)
	sub := &Subscriber{
		ID:       uuid.NewString(),
		Topic:    topic,
		Ch:       make(chan *Envelope, bufferSize),
		cancel:   cancel,
		metadata: metadata,
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.topics[topic]; !ok {
		b.topics[topic] = make(map[string]*Subscriber)
	}
	b.topics[topic][sub.ID] = sub
	activeSubscribers.WithLabelValues(topic).Inc()

	// track lifecycle
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		<-ctx.Done()
		// cleanup when context cancelled
		b.mu.Lock()
		if subs, ok := b.topics[topic]; ok {
			delete(subs, sub.ID)
			if len(subs) == 0 {
				delete(b.topics, topic)
			}
		}
		b.mu.Unlock()
		activeSubscribers.WithLabelValues(topic).Dec()
		close(sub.Ch)
	}()
	return sub, nil
}

// Unsubscribe cancels the subscriber and cleans up.
func (b *Broker) Unsubscribe(sub *Subscriber) {
	if sub == nil {
		return
	}
	sub.cancel()
}

// Publish sends envelope to topic subscribers or direct-to recipient (To).
// Non-blocking: if subscriber channel is full, the envelope will be dropped and logged.
func (b *Broker) Publish(env *Envelope) error {
	if env == nil {
		return errors.New("nil envelope")
	}
	// ensure CreatedAt
	if env.CreatedAt.IsZero() {
		env.CreatedAt = time.Now().UTC()
	}
	messagesPublished.WithLabelValues(env.Topic).Inc()

	// If To set and looks like a subscriber id, attempt direct delivery
	if env.To != "" {
		b.mu.RLock()
		for _, subs := range b.topics {
			if sub, ok := subs[env.To]; ok {
				select {
				case sub.Ch <- env:
					b.mu.RUnlock()
					return nil
				default:
					// channel full; drop
					messagesDropped.WithLabelValues(sub.Topic).Inc()
					b.mu.RUnlock()
					log.Printf("broker: drop direct msg to %s (full)", env.To)
					return nil
				}
			}
		}
		b.mu.RUnlock()
		// fallthrough to topic broadcast if To not found
	}

	// Topic broadcast
	if env.Topic == "" {
		return errors.New("must supply topic or to")
	}
	b.mu.RLock()
	subs, ok := b.topics[env.Topic]
	if !ok || len(subs) == 0 {
		b.mu.RUnlock()
		// no subscribers; safe to return
		return nil
	}

	// deliver to every subscriber non-blocking
	for _, sub := range subs {
		select {
		case sub.Ch <- env:
			// delivered
		default:
			// subscriber channel full -> drop
			messagesDropped.WithLabelValues(env.Topic).Inc()
			log.Printf("broker: drop msg for topic=%s subscriber=%s (full)", env.Topic, sub.ID)
		}
	}
	b.mu.RUnlock()
	return nil
}

// Shutdown stops broker and waits for goroutines to finish.
func (b *Broker) Shutdown(ctx context.Context) {
	b.cancel()
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return
	case <-ctx.Done():
		// forced timeout
		return
	}
}
