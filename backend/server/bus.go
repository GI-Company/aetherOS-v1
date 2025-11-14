package server

import (
	"sync"
	"time"
)

// BusServer: in-memory message broker
type BusServer struct {
	mu           sync.RWMutex
	clients      map[string]*Client
	serverSubs   map[string][]func(*Envelope)
	requests     map[string]chan *Envelope
	metrics      map[string]int64
}

// NewBusServer returns initialized BusServer
func NewBusServer() *BusServer {
	return &BusServer{
		clients:    make(map[string]*Client),
		serverSubs: make(map[string][]func(*Envelope)),
		requests:   make(map[string]chan *Envelope),
		metrics:    make(map[string]int64),
	}
}

// RegisterClient registers a client
func (b *BusServer) RegisterClient(c *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[c.ID] = c
}

// UnregisterClient removes client
func (b *BusServer) UnregisterClient(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if c, ok := b.clients[id]; ok {
		close(c.Send)
		delete(b.clients, id)
	}
}

// Publish sends envelope to server subscribers and clients subscribed to the topic.
func (b *BusServer) Publish(env *Envelope) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	env.Time = time.Now()
	// server-side handlers
	if handlers, ok := b.serverSubs[env.Topic]; ok {
		for _, h := range handlers {
			go h(env)
		}
	}
	// forward to clients (simple prefix matching)
	for _, c := range b.clients {
		// quick check of subscriptions
		for sub := range c.Subscriptions {
			if sub == env.Topic || (len(sub) > 0 && sub[len(sub)-1] == '*' && len(env.Topic) >= len(sub)-1 && env.Topic[:len(sub)-1] == sub[:len(sub)-1]) {
				// non-blocking send
				select {
				case c.Send <- env:
				default:
					// drop if client send buffer full
				}
				break
			}
		}
	}
	b.metrics["published"]++
}

// SubscribeServer registers a server-side handler for a topic
func (b *BusServer) SubscribeServer(topic string, handler func(*Envelope)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.serverSubs[topic] = append(b.serverSubs[topic], handler)
}

// Request sends a request envelope and waits for reply channel with timeout
func (b *BusServer) Request(req *Envelope, timeout time.Duration) (*Envelope, bool) {
	id := time.Now().Format("20060102150405.000000")
	ch := make(chan *Envelope, 1)
	b.mu.Lock()
	b.requests[id] = ch
	b.mu.Unlock()

	// attach request id
	if req.Payload == nil {
		req.Payload = map[string]interface{}{}
	}
	req.Payload["_request_id"] = id
	// publish it
	b.Publish(req)
	// wait
	select {
	case resp := <-ch:
		return resp, true
	case <-time.After(timeout):
		b.mu.Lock()
		delete(b.requests, id)
		b.mu.Unlock()
		return nil, false
	}
}

// ReplyToRequest finds the request channel by id and sends reply
func (b *BusServer) ReplyToRequest(requestID string, reply *Envelope) bool {
	b.mu.Lock()
	ch, ok := b.requests[requestID]
	if ok {
		delete(b.requests, requestID)
	}
	b.mu.Unlock()
	if ok {
		select {
		case ch <- reply:
		default:
		}
		return true
	}
	return false
}
