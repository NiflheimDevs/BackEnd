package repositories

import (
	"context"
	"errors"
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

	query := "SELECT id,owner_id,title,description FROM project WHERE id = $1"

	err := repo.PG.QueryRow(ctx, query, projectID).Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description)

	if err == pgx.ErrNoRows {
		return nil, errors.New("project not found")
	}

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

func (repo *ProjectRepo) CreateProject(userID int, title, description, duration string) int {
	var project_id int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO project (owner_id, title, description, duration) VALUES ($1, $2, $3, $4) RETURNING id"

	err := repo.PG.QueryRow(ctx, query, userID, title, description, duration).Scan(&project_id)

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

func (repo *ProjectRepo) AddProjectTag(projectID, tagID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users_project_tag (project_user_id, tag_id, type) VALUES ($1, $2, 2)"

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

func (repo *ProjectRepo) DeleteProjectTags(projectID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM users_project_tag WHERE project_user_id = $1"

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

func (repo *ProjectRepo) UpdateProject(projectID, UserID int, title, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE project SET title = $1, description = $2 WHERE id = $3 and owner_id = $4"

	result, err := repo.PG.Exec(ctx, query, title, description, projectID, UserID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	if result.RowsAffected() == 0 {
		return errors.New("")
	}

	return nil
}
