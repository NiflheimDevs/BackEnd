package dto

import "time"

type NotifDTO struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
	ReadTime time.Time `json:"read_time"`
	IsRead   bool      `json:"is_read"`
}
