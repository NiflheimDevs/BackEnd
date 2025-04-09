package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PgxTx struct {
	tx pgx.Tx
}

func (t *PgxTx) Exec(ctx context.Context, sql string, args ...any) (any, error) {
	conn, err := t.tx.Exec(ctx, sql, args...)
	return conn.RowsAffected(), err
}

func (t *PgxTx) Query(ctx context.Context, sql string, args ...any) (any, error) {
	return t.tx.Query(ctx, sql, args...)
}

func (t *PgxTx) QueryRow(ctx context.Context, sql string, args ...any) any {
	return t.tx.QueryRow(ctx, sql, args...)
}

func (t *PgxTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *PgxTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
