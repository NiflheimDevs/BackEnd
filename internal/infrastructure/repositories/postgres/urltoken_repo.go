package repositoriesimpl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type UrlTokenRepo struct {
	DB *pgxpool.Pool
}

func NewUrlTokenRepo(db *pgxpool.Pool) *UrlTokenRepo {
	return &UrlTokenRepo{DB: db}
}

func (r *UrlTokenRepo) InsertToken(ctx context.Context, token *models.UrlToken) error {
	query := `
        INSERT INTO token_urls (token, purpose, user_id, invited_email, team_id, sender_id, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.DB.Exec(ctx, query,
		token.Token,
		token.Purpose,
		token.UserID,
		token.InvitedEmail,
		token.TeamID,
		token.SenderID,
		token.ExpiresAt,
	)
	return err
}

func (r *UrlTokenRepo) GetTokenByValue(ctx context.Context, value string) (*models.UrlToken, error) {
	query := `
        SELECT id, token, purpose, user_id, invited_email, team_id, sender_id, expires_at, created_at
        FROM token_urls
        WHERE token = $1
    `
	row := r.DB.QueryRow(ctx, query, value)

	var t models.UrlToken
	err := row.Scan(
		&t.ID,
		&t.Token,
		&t.Purpose,
		&t.UserID,
		&t.InvitedEmail,
		&t.TeamID,
		&t.SenderID,
		&t.ExpiresAt,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}
