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

	var mode int

	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	var level sql.NullInt16

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

	rows, err := gr.PG.Query(ctx, query, careerUserid, mode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []dto.GetTagDto
	for rows.Next() {
		var tag dto.GetTagDto
		err := rows.Scan(&tag.ID, &tag.Name, &level)
		if err != nil {
			return nil, err
		}
		if level.Valid {
			tag.Level = int(level.Int16)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (gr *GeneralRepo) DeleteTagForUserOrCareer(tagid int, careerUserid int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mode int

	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	query := `
        DELETE 
        FROM 
            user_career_tag
        WHERE 
            type = $2 AND career_user_id = $1 AND tag_id = $3`

	_, err := gr.PG.Exec(ctx, query, careerUserid, mode, tagid)
	return err
}

func (gr *GeneralRepo) UpdateTagForUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mode int

	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	query := `
        UPDATE user_career_tag
        SET 
            level = $4
        WHERE 
            type = $2 AND career_user_id = $1 AND tag_id = $3`

	_, err := gr.PG.Exec(ctx, query, careerUserid, mode, tag.ID, tag.Level)
	return err
}

func (gr *GeneralRepo) AddTagToUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mode int

	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	query := `
        INSERT INTO user_career_tag
        (tag_id, career_user_id, type,level)
		VALUES ($1,$2,$3,$4)
			`

	_, err := gr.PG.Exec(ctx, query, tag.ID, careerUserid, mode, tag.Level)
	return err
}

func (gr *GeneralRepo) GetCareersForUser(userid int) ([]models.CareerModel, error) {
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

func (gr *GeneralRepo) CreateCareer(userid int, params *dto.CareerDTO) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var careerid int

	// if params.Website == "" && params.EndDate.IsZero() {

	query := `INSERT INTO career 
	(user_id ,company, start_date,end_date,role,website) 
	VALUES ($1,$2,$3,$4,$5,$6) 
	RETURNING id`

	err := gr.PG.QueryRow(ctx, query, userid, params.Company, params.StartDate, params.EndDate, params.Role, params.Website).Scan(&careerid)

	return careerid, err
}

func (gr *GeneralRepo) UpdateCareer(userid int, params *dto.CareerDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE career 
	SET company = $1 ,start_date = $2 ,end_date = $3 ,role = $4 ,website = $5  
	WHERE user_id = $6 AND id = $7`

	_, err := gr.PG.Exec(ctx, query, params.Company, params.StartDate, params.EndDate, params.Role, params.Website, userid, params.ID)

	return err
}

func (gr *GeneralRepo) DeleteTagForCareerOrUserByID(tagid int, careerUserid int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mode int
	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	query := `DELETE
	FROM users_career_tag
	WHERE tag_id = $1 AND career_user_tag = $2 AND type = $3`

	_, err := gr.PG.Exec(ctx, query, tagid, careerUserid, mode)
	return err
}

func (gr *GeneralRepo) DeleteCareer(userid int, careerid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE
	FROM users_career
	WHERE user_id = $1 AND career_id = $2`

	_, err := gr.PG.Exec(ctx, query, userid, careerid)
	return err
}
