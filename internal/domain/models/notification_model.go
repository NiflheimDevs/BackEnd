package models

import "time"

type NotifModel struct {
	ID       int
	UserID   int
	Content  string
	SendTime time.Time
	ReadTime time.Time
	IsRead   bool
}
