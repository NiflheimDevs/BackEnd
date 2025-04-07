package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/models"
)

type ProjectRepo struct {
	PG *pgxpool.Pool
}

func NewProjectRepo(
	PG *pgxpool.Pool,
) *ProjectRepo {
	return &ProjectRepo{
		PG: PG,
	}
}

func (repo *ProjectRepo) GetProject(projectID int) (*models.ProjectModel, error) {
	var project models.ProjectModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,owner_id,title,description,label,duration FROM project WHERE id = $1"

	var duration time.Time
	err := repo.PG.QueryRow(ctx, query, projectID).Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &project.Label, &duration)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	project.Duration = duration.Format("2006-01-02 15:04:05")

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	return &project, nil
}

func (repo *ProjectRepo) GetUserProject(userID, offset, limit int) []models.ProjectModel {
	var projects []models.ProjectModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,owner_id,title,description,label,duration FROM project WHERE owner_id = $1 ORDER BY id OFFSET $2 LIMIT $3"

	result, err := repo.PG.Query(ctx, query, userID, offset, limit)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	defer result.Close()

	for result.Next() {
		var project models.ProjectModel
		var duration time.Time
		if err := result.Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &project.Label, &duration); err != nil {
			panic(exceptions.Exception{
				Tag: enums.INTERNAL_ERROR,
				Errors: []enums.SpecificError{
					enums.DATABASE_ERROR,
				},
			})
		}

		project.Duration = duration.Format("2006-01-02 15:04:05")
		tag := repo.GetProjectTag(project.ID)
		project.Tags = tag
		projects = append(projects, project)
	}

	return projects
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, tx pgx.Tx, userID, label int, title, description, duration string) int {
	var project_id int

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO project (owner_id, title, description, label, duration, created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	err := tx.QueryRow(ctx, query, userID, title, description, label, duration, now, now).Scan(&project_id)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	return project_id
}

func (repo *ProjectRepo) GetProjectTag(projectID int) []models.TagModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tags []models.TagModel

	query := "SELECT t.id, t.name FROM project_tag pt JOIN tag t ON pt.tag_id = t.id WHERE pt.project_id = $1"

	rows, err := repo.PG.Query(ctx, query, projectID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.TagModel
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			panic(exceptions.Exception{
				Tag: enums.INTERNAL_ERROR,
				Errors: []enums.SpecificError{
					enums.DATABASE_ERROR,
				},
			})
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	return tags
}

func (repo *ProjectRepo) AddProjectTag(ctx context.Context, tx pgx.Tx, projectID, tagID int) {
	query := "INSERT INTO project_tag (project_id, tag_id) VALUES ($1, $2)"

	_, err := tx.Exec(ctx, query, projectID, tagID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) DeleteProjectTags(ctx context.Context, tx pgx.Tx, projectID, tagID int) {
	query := "DELETE FROM project_tag WHERE project_id = $1 AND tag_id = $2"

	_, err := tx.Exec(ctx, query, projectID, tagID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) UpdateProject(ctx context.Context, tx pgx.Tx, projectID, UserID, label int, title, description string) {
	now := time.Now().Format("2006-01-02 15:04:05")

	query := "UPDATE project SET title = $1, description = $2, label=$3, updated_time=$4 WHERE id = $5 and owner_id = $6"

	_, err := tx.Exec(ctx, query, title, description, label, now, projectID, UserID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) DeleteProject(ctx context.Context, tx pgx.Tx, projectID int) {
	query := "DELETE FROM project WHERE id = $1"
	_, err := tx.Exec(ctx, query, projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}
