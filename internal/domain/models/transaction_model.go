package models

import "time"

type TransactionModel struct {
	ID          int
	FromUser    int
	ToUser      int
	Amount      int64
	Description string
	Date        time.Time
}
