
// =================================
// backend/server/types.go
// =================================
package server

type Message struct {
	Topic   string                 `json:"topic"`
	Payload map[string]interface{} `json:"payload"`
	Token   string                 `json:"token,omitempty"`
	ReplyTo string                 `json:"replyTo,omitempty"`
	Source  string                 `json:"source,omitempty"`
	Dest    string                 `json:"dest,omitempty"`
}
