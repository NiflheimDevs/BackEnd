package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/models"
)

type GeneralRepo struct {
	PG *pgxpool.Pool
}

func NewGeneralRepo(pg *pgxpool.Pool) *GeneralRepo {
	return &GeneralRepo{
		PG: pg,
	}
}

func (repo *GeneralRepo) GetTags() ([]models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM tag`
	result, err := repo.PG.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	var tags []models.TagModel
	for result.Next() {
		var tag models.TagModel
		err = result.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (gr *GeneralRepo) GetTagsForUserOrCareer(careerUserid int, isForUser bool) ([]models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mod int

	if isForUser {
		mod = 0
	} else {
		mod = 1
	}

	query := `
        SELECT 
            t.id, 
            t.name
        FROM 
            user_career_tag AS uct
        JOIN 
            tag AS t 
        ON 
            t.id = uct.tag_id
        WHERE 
            uct.type = $2 
            AND uct.career_user_id = $1`

	rows, err := gr.PG.Query(ctx, query, careerUserid, mod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.TagModel
	for rows.Next() {
		var tag models.TagModel
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (gr *GeneralRepo) GetCareerForUser(userid int) ([]models.CareerModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.id,
			c.company,
			c.start_date,
			c.end_date,
			c.role,
			c.website,
		FROM 
			career AS c
		WHERE
			c.user_id = $1`

	rows, err := gr.PG.Query(ctx, query, userid)
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

func (gr *GeneralRepo) PostCareer() {

}
