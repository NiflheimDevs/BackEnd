package pkg

import "github.com/google/uuid"

func NewUniqueID() string {
	return uuid.New().String()
}
