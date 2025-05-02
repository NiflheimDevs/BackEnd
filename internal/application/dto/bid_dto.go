package dto

import "time"

type BidInfo struct {
	BidID        int       `json:"bid_id" validate:""`
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

type PublicProjectBidInfo struct {
	BidID        int    `json:"bid_id"`
	Title        string `json:"title"`
	Total        int64  `json:"total"`
	ExpectedTime string `json:"expected_time"`
	ProfilePic   string `json:"profile_pic"`
}

type PrivateProjectBidInfo struct {
	BidID        int            `json:"bid_id"`
	UserInfo     UserProfileDTO `json:"user_info"`
	PrePayment   int64          `json:"pre_payment"`
	Total        int64          `json:"total"`
	ExpectedTime string         `json:"expected_time"`
}
