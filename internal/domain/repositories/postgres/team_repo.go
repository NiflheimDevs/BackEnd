package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TeamRepo interface {
	CreateOneManTeam(ctx context.Context, tx transaction.Tx, userid int) int64
	AddMember(userid int, teamid int64, position string) error
	AddMemberWithTx(ctx context.Context, tx transaction.Tx, userid int, teamid int64, position string) error
}
