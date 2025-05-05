package dto

type BidInfo struct {
	BidID        int    `json:"bid_id" validate:""`
	UserID       int    `json:"user_id" validate:""`
	TeamID       int64  `json:"team_id" validate:"required"`
	ProjectID    int    `json:"project_id" validate:"required"`
	PP           int64  `json:"pre_payment" validate:"required"`
	Total        int64  `json:"total" validate:"required"`
	Description  string `json:"description" validate:""`
	ExpectedTime int    `json:"expected_time" validate:"required"`
}

type PutBid struct {
	BidID int `json:"bid_id"`
}

type PublicProjectBidInfo struct {
	BidID        int                  `json:"bid_id"`
	TeamInfo     *GetInternalTeamInfo `json:"team_info"`
	Total        int64                `json:"total"`
	ExpectedTime int                  `json:"expected_time"`
}

type PrivateProjectBidInfo struct {
	BidID        int                  `json:"bid_id"`
	Type         int                  `json:"type"`
	TeamInfo     *GetInternalTeamInfo `json:"team_info"`
	PrePayment   int64                `json:"pre_payment"`
	Total        int64                `json:"total"`
	ExpectedTime int                  `json:"expected_time"`
}
