package models

import "time"

type BidModel struct {
	ID           int
	TeamID       int64
	ProjectID    int
	PrePayment   int64
	Total        int64
	Description  string
	ExpectedTime time.Time
	CreatedTime  time.Time
}
