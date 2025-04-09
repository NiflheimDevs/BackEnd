package repositoriesimpl

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type ProjectRepo struct {
	PG      *pgxpool.Pool
	TagRepo TagRepo
}

func NewProjectRepo(
	PG *pgxpool.Pool,
	tagRepo TagRepo,
) *ProjectRepo {
	return &ProjectRepo{
		PG:      PG,
		TagRepo: tagRepo,
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
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return &project, nil
}

func (repo *ProjectRepo) GetUserProject(userID, offset, limit int) []models.ProjectModel {
	var projects []models.ProjectModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,owner_id,title,description,label,duration FROM project WHERE owner_id = $1 OFFSET $2 LIMIT $3"

	result, err := repo.PG.Query(ctx, query, userID, offset, limit)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	defer result.Close()

	for result.Next() {
		var project models.ProjectModel
		var duration time.Time
		if err := result.Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &project.Label, &duration); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}

		project.Duration = duration.Format("2006-01-02 15:04:05")
		tag := repo.TagRepo.GetProjectTag(project.ID)
		project.Tags = tag
		projects = append(projects, project)
	}

	return projects
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, tx transaction.Tx, userID, label int, title, description, duration string) int {
	var project_id int

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO project (owner_id, title, description, label, duration, created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	row := tx.QueryRow(ctx, query, userID, title, description, label, duration, now, now).(pgx.Row)
	err := row.Scan(&project_id)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return project_id
}

func (repo *ProjectRepo) UpdateProject(projectID, UserID int, title, description string) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "UPDATE project SET title = $1, description = $2, label=$3, updated_time=$4 WHERE id = $5 and owner_id = $6"

	_, err := repo.PG.Exec(ctx, query, title, description, now, projectID, UserID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) DeleteProject(ctx context.Context, tx transaction.Tx, projectID int) {
	query := "DELETE FROM project WHERE id = $1"
	_, err := tx.Exec(ctx, query, projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}
