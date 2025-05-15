package websocket

import (
	"encoding/json"
	"time"
)

const (
	MessageTypeChat         = "chat"
	MessageTypeNotification = "notification"
)

type Message struct {
	MessageID uint            `json:"id"`
	Type      string          `json:"type"`
	RoomID    uint            `json:"room_id,omitempty"`
	SenderID  uint            `json:"sender_id,omitempty"`
	Content   json.RawMessage `json:"content"`
	Timestamp time.Time       `json:"timestamp"`
	Client    *Client         `json:"-"`
}
