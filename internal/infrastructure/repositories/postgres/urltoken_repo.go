package repositoriesimpl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type UrlTokenRepo struct {
	DB *pgxpool.Pool
}

func NewUrlTokenRepo(db *pgxpool.Pool) *UrlTokenRepo {
	return &UrlTokenRepo{DB: db}
}

func (r *UrlTokenRepo) InsertToken(token *models.UrlToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO url_token (token, purpose, user_id, invited_email, team_id, sender_id, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.DB.Exec(ctx, query,
		token.Token,
		token.Purpose,
		token.UserID,
		// token.InvitedEmail,
		token.TeamID,
		token.SenderID,
		token.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
	}

	return err
}

func (r *UrlTokenRepo) DeleteToken(value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM url_token WHERE token = $1;
	`
	_, err := r.DB.Exec(ctx, query, value)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
	}

	return err
}

func (r *UrlTokenRepo) DeleteTokenByTeamAndPurposeAndUser(teamid int64, userid int, Purpose enums.UrlTokenPurpose) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM url_token WHERE team_id = $1 AND purpose = $2 AND user_id = $3;
	`
	_, err := r.DB.Exec(ctx, query, teamid, Purpose, userid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
	}

	return err
}

func (r *UrlTokenRepo) DeleteTokenByPurposeAndUser(userid int, Purpose enums.UrlTokenPurpose) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM url_token WHERE purpose = $1 AND user_id = $2;
	`
	_, err := r.DB.Exec(ctx, query, Purpose, userid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
	}

	return err
}

func (r *UrlTokenRepo) GetTokenByValue(value string) (*models.UrlToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT id, token, purpose, user_id, invited_email, team_id, sender_id, expires_at, created_at
        FROM url_token
        WHERE token = $1
    `
	row := r.DB.QueryRow(ctx, query, value)

	var exp sql.NullTime

	var t models.UrlToken
	err := row.Scan(
		&t.ID,
		&t.Token,
		&t.Purpose,
		&t.UserID,
		// &t.InvitedEmail,
		&t.TeamID,
		&t.SenderID,
		&exp,
		&t.CreatedAt,
	)

	if exp.Valid {
		t.ExpiresAt = exp.Time
	}
	// new pattern
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil, err
	}

	return &t, nil
}

func (r *UrlTokenRepo) GetTokenByValueAndPurpose(value string, purpose enums.UrlTokenPurpose) (*models.UrlToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT id, token, purpose, user_id, invited_email, team_id, sender_id, expires_at, created_at
        FROM url_token
        WHERE token = $1 AND purpose = $2
    `

	row := r.DB.QueryRow(ctx, query, value, purpose)

	var exp sql.NullTime

	var t models.UrlToken
	err := row.Scan(
		&t.ID,
		&t.Token,
		&t.Purpose,
		&t.UserID,
		// &t.InvitedEmail,
		&t.TeamID,
		&t.SenderID,
		&exp,
		&t.CreatedAt,
	)

	if exp.Valid {
		t.ExpiresAt = exp.Time
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil, err
	}

	return &t, nil
}

func (r *UrlTokenRepo) GetTokenByTeamAndPurposeAndUser(teamid int64, userid int, Purpose enums.UrlTokenPurpose) (*models.UrlToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT id, token, purpose, user_id, invited_email, team_id, sender_id, expires_at, created_at
        FROM url_token
        WHERE team_id = $1 AND purpose = $2 AND user_id = $3;
    `

	row := r.DB.QueryRow(ctx, query, teamid, Purpose, userid)

	var exp sql.NullTime

	var t models.UrlToken
	err := row.Scan(
		&t.ID,
		&t.Token,
		&t.Purpose,
		&t.UserID,
		// &t.InvitedEmail,
		&t.TeamID,
		&t.SenderID,
		&exp,
		&t.CreatedAt,
	)

	if exp.Valid {
		t.ExpiresAt = exp.Time
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil, err
	}

	return &t, nil
}

func (r *UrlTokenRepo) GetTokenByPurposeAndUser(userid int, Purpose enums.UrlTokenPurpose) (*models.UrlToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        SELECT id, token, purpose, user_id, invited_email, team_id, sender_id, expires_at, created_at
        FROM url_token
        WHERE purpose = $1 AND user_id = $2;
    `

	row := r.DB.QueryRow(ctx, query, Purpose, userid)

	var exp sql.NullTime

	var t models.UrlToken
	err := row.Scan(
		&t.ID,
		&t.Token,
		&t.Purpose,
		&t.UserID,
		// &t.InvitedEmail,
		&t.TeamID,
		&t.SenderID,
		&exp,
		&t.CreatedAt,
	)

	if exp.Valid {
		t.ExpiresAt = exp.Time
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil, err
	}

	return &t, nil
}
