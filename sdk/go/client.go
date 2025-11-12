
package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// Client is a client for the Aether broker.
type Client struct {
	BrokerURL string
	Token     string
	client    *http.Client
}

// NewClient creates a new Aether client.
func NewClient(brokerURL, token string) (*Client, error) {
	return &Client{
		BrokerURL: brokerURL,
		Token:     token,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Envelope is a message envelope.
type Envelope struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

// UnmarshalPayload unmarshals the envelope's payload into the given value.
func (e *Envelope) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}

// Publish publishes a message to a topic.
func (c *Client) Publish(ctx context.Context, topic string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	env := &Envelope{
		Topic:   topic,
		Payload: payloadBytes,
	}

	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("error marshaling envelope: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BrokerURL+"/v1/bus/publish", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Subscribe subscribes to a topic.
func (c *Client) Subscribe(ctx context.Context, topic string) (<-chan *Envelope, error) {
	u, err := url.Parse(c.BrokerURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing broker URL: %w", err)
	}

	// Use wss for https, ws for http
	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}

	u.Scheme = scheme
	u.Path = "/v1/bus/subscribe"
	q := u.Query()
	q.Set("topic", topic)
	u.RawQuery = q.Encode()

	h := http.Header{}
	h.Set("Authorization", "Bearer "+c.Token)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket: %w", err)
	}

	msgs := make(chan *Envelope)
	go func() {
		defer close(msgs)
		defer conn.Close()
		for {
			var env Envelope
			if err := conn.ReadJSON(&env); err != nil {
				return
			}
			select {
			case msgs <- &env:
			case <-ctx.Done():
				return
			}
		}
	}()

	return msgs, nil
}

// NewJWT creates a new JWT for the given user ID.
func NewJWT(userID string, ttl time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// In a real application, you would use a secret from a secure source.
	return token.SignedString([]byte("aether-secret"))
}
