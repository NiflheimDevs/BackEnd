package models

import "time"

type TeamModel struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Type        int       `json:"type"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
}
