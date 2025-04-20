package services

import (
	"context"

	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TeamService interface {
	BehindCurtainTeam(ctx context.Context, tx transaction.Tx, userid int)
}
