package repositoriesimpl

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type CareerRepo struct {
	PG *pgxpool.Pool
}

func NewCareerRepo(pg *pgxpool.Pool) *CareerRepo {
	return &CareerRepo{
		PG: pg,
	}
}

func (cr *CareerRepo) GetCareersForUser(userid int) ([]models.CareerModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.id,
			c.company,
			c.start_date,
			c.end_date,
			c.role,
			c.website
		FROM 
			career AS c
		WHERE
			c.user_id = $1`

	rows, err := cr.PG.Query(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var careers []models.CareerModel
	var website sql.NullString
	var endDate sql.NullTime
	for rows.Next() {
		var career models.CareerModel
		err := rows.Scan(&career.ID, &career.Company, &career.StartDate, &endDate, &career.Role, &website)
		if err != nil {
			return nil, err
		}
		if website.Valid {
			career.Website = website.String
		}
		if endDate.Valid {
			career.EndDate = endDate.Time
		}

		careers = append(careers, career)
	}

	return careers, nil
}

func (cr *CareerRepo) CreateCareer(userid int, params *dto.CareerDTO) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var careerid int

	query := `INSERT INTO career 
	(user_id ,company, start_date,end_date,role,website) 
	VALUES ($1,$2,$3,$4,$5,$6) 
	RETURNING id`

	err := cr.PG.QueryRow(ctx, query, userid, params.Company, params.StartDate, params.EndDate, params.Role, params.Website).Scan(&careerid)

	return careerid, err
}

func (cr *CareerRepo) UpdateCareer(userid int, params *dto.CareerDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE career 
	SET company = $1 ,start_date = $2 ,end_date = $3 ,role = $4 ,website = $5  
	WHERE user_id = $6 AND id = $7`

	_, err := cr.PG.Exec(ctx, query, params.Company, params.StartDate, params.EndDate, params.Role, params.Website, userid, params.ID)

	return err
}

func (cr *CareerRepo) DeleteCareer(userid int, careerid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE
	FROM users_career
	WHERE user_id = $1 AND career_id = $2`

	_, err := cr.PG.Exec(ctx, query, userid, careerid)
	return err
}
