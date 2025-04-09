package repositoriesimpl

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type TagRepo struct {
	PG *pgxpool.Pool
}

func NewTagRepo(pg *pgxpool.Pool) *TagRepo {
	return &TagRepo{
		PG: pg,
	}
}

func (tr *TagRepo) GetTags() ([]models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM tag`
	result, err := tr.PG.Query(ctx, query)

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

func (tr *TagRepo) GetTagsForUserOrCareer(careerUserid int, isForUser bool) ([]dto.GetTagDto, error) {
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
            users_career_tag AS uct
        JOIN 
            tag AS t 
        ON 
            t.id = uct.tag_id
        WHERE 
            uct.type = $2 
            AND uct.career_user_id = $1`

	rows, err := tr.PG.Query(ctx, query, careerUserid, mode)
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

func (tr *TagRepo) DeleteTagForUserOrCareer(tagid int, careerUserid int, isForUser bool) error {
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
            users_career_tag
        WHERE 
            type = $2 AND career_user_id = $1 AND tag_id = $3`

	_, err := tr.PG.Exec(ctx, query, careerUserid, mode, tagid)
	return err
}

func (tr *TagRepo) UpdateTagForUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mode int

	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	query := `
        UPDATE users_career_tag
        SET 
            level = $4
        WHERE 
            type = $2 AND career_user_id = $1 AND tag_id = $3`

	_, err := tr.PG.Exec(ctx, query, careerUserid, mode, tag.ID, tag.Level)
	return err
}

func (tr *TagRepo) AddTagToUserOrCareer(tag *dto.RecieveTagDTO, careerUserid int, isForUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mode int

	if isForUser {
		mode = 0
	} else {
		mode = 1
	}

	query := `
        INSERT INTO users_career_tag
        (tag_id, career_user_id, type,level)
		VALUES ($1,$2,$3,$4)
			`

	_, err := tr.PG.Exec(ctx, query, tag.ID, careerUserid, mode, tag.Level)
	return err
}

func (tr *TagRepo) DeleteTagForCareerOrUserByID(tagid int, careerUserid int, isForUser bool) error {
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
	WHERE tag_id = $1 AND career_user_id = $2 AND type = $3`

	_, err := tr.PG.Exec(ctx, query, tagid, careerUserid, mode)
	return err
}

func (tr *TagRepo) GetProjectTag(projectID int) []models.TagModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tags []models.TagModel

	query := "SELECT t.id, t.name FROM project_tag pt JOIN tag t ON pt.tag_id = t.id WHERE pt.project_id = $1"

	rows, err := tr.PG.Query(ctx, query, projectID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.TagModel
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return tags
}

func (tr *TagRepo) AddProjectTagWithTx(ctx context.Context, tx transaction.Tx, projectID, tagID int) error {
	query := "INSERT INTO project_tag (project_id, tag_id) VALUES ($1, $2)"

	_, err := tx.Exec(ctx, query, projectID, tagID)

	return err
}

func (tr *TagRepo) DeleteProjectTagWithTx(ctx context.Context, tx transaction.Tx, projectID, tagID int) error {
	query := "DELETE FROM project_tag WHERE project_id = $1 AND tag_id = $2"

	_, err := tx.Exec(ctx, query, projectID, tagID)

	return err
}

func (tr *TagRepo) AddProjectTag(projectID, tagID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO project_tag (project_id, tag_id) VALUES ($1, $2)"

	_, err := tr.PG.Exec(ctx, query, projectID, tagID)

	return err
}

func (tr *TagRepo) DeleteProjectTag(projectID, tagID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM project_tag WHERE project_id = $1 AND tag_id = $2"

	_, err := tr.PG.Exec(ctx, query, projectID, tagID)

	return err
}
