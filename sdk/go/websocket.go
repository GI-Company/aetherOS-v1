package aethersdk

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient represents a bidirectional websocket connection to the broker.
type WSClient struct {
	conn    *websocket.Conn
	sendMu  sync.Mutex       // guards writes to conn
	recvCh  chan *Envelope   // channel for incoming envelopes
	done    chan struct{}    // closed when connection is closed
	closed  bool
	closeMu sync.Mutex
}

// Subscribe now returns a WSClient to both receive and send envelopes.
func (c *Client) Subscribe(ctx context.Context, topic string) (*WSClient, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	// convert scheme
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}
	u.Path = "/v1/bus/subscribe"
	q := u.Query()
	q.Set("topic", topic)
	q.Set("sid", "sdk-go-client")
	u.RawQuery = q.Encode()

	header := make(map[string][]string)
	header["Authorization"] = []string{"Bearer " + c.Token}

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.DialContext(ctx, u.String(), header)
	if err != nil {
		return nil, err
	}

	ws := &WSClient{
		conn:   conn,
		recvCh: make(chan *Envelope, 64),
		done:   make(chan struct{}),
	}

	// reader goroutine
	go func() {
		defer func() {
			ws.closeInternal()
		}()
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("aether-sdk: ws read error: %v", err)
				return
			}
			var env Envelope
			if err := json.Unmarshal(data, &env); err != nil {
				log.Printf("aether-sdk: invalid envelope: %v", err)
				continue
			}
			select {
			case ws.recvCh <- &env:
			default:
				// recv channel full -> drop oldest (non-blocking policy)
				select {
				case <-ws.recvCh:
				default:
				}
				ws.recvCh <- &env
			}
		}
	}()

	// writer goroutine is not required here (Send uses conn directly under mutex)
	return ws, nil
}

// Receive returns the channel of incoming envelopes.
func (w *WSClient) Receive() <-chan *Envelope {
	return w.recvCh
}

// Send sends an envelope on the same websocket connection (thread-safe).
func (w *WSClient) Send(ctx context.Context, env *Envelope) error {
	if env == nil {
		return errors.New("env required")
	}
	w.closeMu.Lock()
	closed := w.closed
	w.closeMu.Unlock()
	if closed {
		return errors.New("ws client closed")
	}
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	// use mutex to ensure single writer
	w.sendMu.Lock()
	defer w.sendMu.Unlock()

	// respect context deadline
	done := make(chan error, 1)
	go func() {
		err := w.conn.WriteMessage(websocket.TextMessage, data)
		done <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil {
			return err
		}
		return nil
	}
}

// Close gracefully closes the websocket client.
func (w *WSClient) Close() {
	w.closeInternal()
}

// internal close with idempotence
func (w *WSClient) closeInternal() {
	w.closeMu.Lock()
	defer w.closeMu.Unlock()
	if w.closed {
		return
	}
	w.closed = true
	_ = w.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
	_ = w.conn.Close()
	close(w.done)
	// drain/close recvCh
	close(w.recvCh)
}
