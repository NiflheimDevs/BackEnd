package models

type BidModel struct {
	ID           int
	TeamID       int
	ProjectID    int
	Value        int64
	ExpectedTime string
	CreatedTime  string
}
