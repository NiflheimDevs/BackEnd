package models

type BidModel struct {
	ID           int
	TeamID       int
	ProjectID    int
	PrePayment   int64
	Total        int64
	Description  string
	ExpectedTime string
	CreatedTime  string
}
