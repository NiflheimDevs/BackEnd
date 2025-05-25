package models

import "time"

type MessageModel struct {
	MessageID int
	RoomID    int
	SenderID  int
	SendTime  time.Time
	EditTime  time.Time
	Content   string
}
