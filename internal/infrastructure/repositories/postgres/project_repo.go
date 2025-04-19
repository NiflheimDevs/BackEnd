package repositoriesimpl

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type ProjectRepo struct {
	PG      *pgxpool.Pool
	TagRepo repositories.TagRepo
}

func NewProjectRepo(
	PG *pgxpool.Pool,
	tagRepo repositories.TagRepo,
) *ProjectRepo {
	return &ProjectRepo{
		PG:      PG,
		TagRepo: tagRepo,
	}
}

func (repo *ProjectRepo) LandingProps() []dto.ProjectLanding {
	var projects []dto.ProjectLanding

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT P.id,P.title,P.description,L.name From project P Join label L on P.label = L.id ORDER BY P.label DESC ,P.created_time DESC LIMIT 4"

	result, err := repo.PG.Query(ctx, query)

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
		var project dto.ProjectLanding
		if err := result.Scan(&project.ProjectID, &project.Title, &project.Description, &project.Label); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		projects = append(projects, project)
	}
	return projects
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

	query := "SELECT id,owner_id,title,description,label,duration FROM project WHERE owner_id = $1 ORDER BY id OFFSET $2 LIMIT $3"

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

	query := "UPDATE project SET title = $1, description = $2, updated_time=$3 WHERE id = $4 and owner_id = $5"

	_, err := repo.PG.Exec(ctx, query, title, description, now, projectID, UserID)

	if err != nil {
		log.Println(err)
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
