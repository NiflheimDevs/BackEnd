package models

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/enums"
)

type UrlToken struct {
	ID           int
	Token        string
	Purpose      enums.UrlTokenPurpose
	UserID       *int
	InvitedEmail string
	TeamID       *int
	SenderID     *int
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
