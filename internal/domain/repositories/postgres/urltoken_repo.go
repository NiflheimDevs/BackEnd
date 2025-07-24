package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/domain/models"
)

type UrlTokenRepo interface {
	InsertToken(ctx context.Context, token *models.UrlToken) error
	GetTokenByValue(ctx context.Context, value string) (*models.UrlToken, error)
}
