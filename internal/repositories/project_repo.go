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

	query := "SELECT id,owner_id,title,description,duration FROM project WHERE id = $1"

	var duration time.Time
	err := repo.PG.QueryRow(ctx, query, projectID).Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &duration)

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

	query := "SELECT id,owner_id,title,description,duration FROM project WHERE owner_id = $1 OFFSET $2 LIMIT $3"

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
		if err := result.Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &duration); err != nil {
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

func (repo *ProjectRepo) CreateProject(userID int, title, description, duration string) int {
	var project_id int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO project (owner_id, title, description, duration, created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	err := repo.PG.QueryRow(ctx, query, userID, title, description, duration, now, now).Scan(&project_id)

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

func (repo *ProjectRepo) AddProjectTag(projectID, tagID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO project_tag (project_id, tag_id) VALUES ($1, $2)"

	_, err := repo.PG.Exec(ctx, query, projectID, tagID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) DeleteProjectTags(projectID, tagID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM project_tag WHERE project_id = $1 AND tag_id = $2"

	_, err := repo.PG.Exec(ctx, query, projectID, tagID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) UpdateProject(projectID, UserID int, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "UPDATE project SET title = $1, description = $2, updated_time=$3 WHERE id = $4 and owner_id = $5"

	_, err := repo.PG.Exec(ctx, query, title, description, now, projectID, UserID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) DeleteProject(projectID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM project WHERE id = $1"
	_, err := repo.PG.Exec(ctx, query, projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}
