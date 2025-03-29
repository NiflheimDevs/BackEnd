package dto

type UserTransactionDTO struct {
	Date        string `json:"date"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}
