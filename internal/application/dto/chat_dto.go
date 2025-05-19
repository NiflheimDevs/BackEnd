package dto

import "time"

type RoomInfo struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	RoomID    int    `json:"room_id"`
}

type Message struct {
	ID       int
	Content  string
	SenderID int
	SendTime time.Time
	Type     int
}
