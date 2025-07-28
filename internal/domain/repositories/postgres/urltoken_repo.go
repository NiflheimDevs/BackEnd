package repositories

import (
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type UrlTokenRepo interface {
	InsertToken(token *models.UrlToken) error
	DeleteToken(value string) error
	GetTokenByValue(value string) (*models.UrlToken, error)
	GetTokenByValueAndPurpose(value string, purpose enums.UrlTokenPurpose) (*models.UrlToken, error)
	GetTokenByPurposeAndUser(userid int, Purpose enums.UrlTokenPurpose) (*models.UrlToken, error)
	GetTokenByTeamAndPurposeAndUser(teamid int64, userid int, Purpose enums.UrlTokenPurpose) (*models.UrlToken, error)
	DeleteTokenByPurposeAndUser(userid int, Purpose enums.UrlTokenPurpose) error
	DeleteTokenByTeamAndPurposeAndUser(teamid int64, userid int, Purpose enums.UrlTokenPurpose) error
}
