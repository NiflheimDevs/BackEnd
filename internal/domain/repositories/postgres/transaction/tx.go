package transaction

import (
	"context"
)

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (any, error)
	Query(ctx context.Context, sql string, args ...any) (any, error)
	QueryRow(ctx context.Context, sql string, args ...any) any

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
