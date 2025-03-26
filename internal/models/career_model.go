package models

import "time"

type CareerModel struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Company   string    `json:"company"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Role      string    `json:"role"`
	Website   string    `json:"website"`
}
