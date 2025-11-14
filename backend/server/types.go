package server

import (
	"time"

	"github.com/gorilla/websocket"
)

// Envelope is the canonical message carried on the message bus
type Envelope struct {
	Topic   string                 `json:"topic"`
	From    string                 `json:"from,omitempty"`
	To      string                 `json:"to,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
	Time    time.Time              `json:"time"`
}

// MessageRequest for synchronous request/response
type MessageRequest struct {
	ID      string
	Topic   string
	Payload map[string]interface{}
	ReplyCh chan *Envelope
}

// Client represents a websocket client connection
type Client struct {
	ID            string
	Conn          *websocket.Conn
	Send          chan *Envelope
	Subscriptions map[string]bool
}

// WindowState and App snapshot types
type WindowState struct {
	ID    string `json:"id"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	W     int    `json:"w"`
	H     int    `json:"h"`
	Z     int64  `json:"z"`
	State string `json:"state"`
}

// AppSnapshot stored per app
type AppSnapshot struct {
	AppID     string                 `json:"appId"`
	Windows   []WindowState          `json:"windows"`
	AppState  map[string]interface{} `json:"appState"`
	Dirty     bool                   `json:"dirty"`
	LastSaved time.Time              `json:"lastSaved"`
	Version   int                    `json:"version"`
}
