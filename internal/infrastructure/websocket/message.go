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
	MessageID int             `json:"id"`
	Type      string          `json:"type"`
	RoomID    int             `json:"room_id"`
	SenderID  int             `json:"sender_id"`
	Content   json.RawMessage `json:"content"`
	Timestamp time.Time       `json:"timestamp"`
	Client    *Client         `json:"-"`
}
