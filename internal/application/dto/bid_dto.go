package dto

import "time"

type BidInfo struct {
	UserID       int       `json:"user_id" validate:""`
	TeamID       int       `json:"team_id" validate:"required"`
	ProjectID    int       `json:"project_id" validate:"required"`
	PP           int64     `json:"pre_payment" validate:"required"`
	Total        int64     `json:"total" validate:"required"`
	Description  string    `json:"description" validate:""`
	ExpectedTime time.Time `json:"expected_time" validate:"required"`
}

type PutBid struct {
	BidID int `json:"bid_id"`
}

type ProjectBidInfo struct {
	BidID int   `json:"bid_id"`
	Value int64 `json:"value"`
}
