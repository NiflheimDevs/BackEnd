package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/dto"
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

func (gr *GeneralRepo) GetTagsForUserOrCareer(careerUserid int, isForUser bool) ([]dto.GetTagDto, error) {
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
            t.name,
			uct.level
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

	var tags []dto.GetTagDto
	for rows.Next() {
		var tag dto.GetTagDto
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

func (gr *GeneralRepo) PostCareer(userid int, params *dto.PostCareerDTO) (*models.CareerModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// var query string
	var careerid int
	// var err error

	// if params.Website == "" && params.EndDate.IsZero() {

	// 	query = `INSERT INTO career (user_id ,company, start_date,role) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	// 	err = gr.PG.QueryRow(ctx, query, userid, params.Company, params.StartDate, params.Role).Scan(&careerid)
	// }

	query := `INSERT INTO career (user_id ,company, start_date,end_date,role,website ) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	err := gr.PG.QueryRow(ctx, query, userid, params.Company, params.StartDate, params.EndDate, params.Role, params.Website).Scan(&careerid)
	response := models.CareerModel{
		ID:        careerid,
		Company:   params.Company,
		UserID:    userid,
		StartDate: params.StartDate,
		Role:      params.Role,
		Website:   params.Website,
		EndDate:   params.EndDate,
	}

	return &response, err
}

func (gr *GeneralRepo) AddTagToUserOrCareer(careerUserid int, tagid int, level int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mod int
	if isForUser {
		mod = 0
	} else {
		mod = 1
	}

	query := `INSERT INTO users_career_tag 
	(tag_id, career_user_id, type, level)
	VALUES ($1,$2,$3,$4)`

	_, err := gr.PG.Exec(ctx, query, tagid, careerUserid, mod, level)

	return err
}
