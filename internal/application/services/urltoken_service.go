package services

import (
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type UrlTokenService interface {
	GenerateTeamInviteToken(userid, senderID int, teamid int64) (string, error)
	GenerateEmailVerificationToken(userid int) (string, error)
	GetToken(tokenStr string) (*models.UrlToken, error)
	GetTokenWithPurpose(tokenStr string, expectedPurpose enums.UrlTokenPurpose) (*models.UrlToken, error)
	DeleteToken(token string) error
}
