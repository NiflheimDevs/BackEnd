package dto

type PutBid struct {
	BidID int `json:"bid_id"`
}

type ProjectBidInfo struct {
	BidID int   `json:"bid_id"`
	Value int64 `json:"value"`
}
