package repositoriesimpl

import (
	"context"
	"database/sql"
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
	PG          *pgxpool.Pool
	TagRepo     repositories.TagRepo
	CommentRepo repositories.CommentRepo
}

func NewProjectRepo(
	PG *pgxpool.Pool,
	tagRepo repositories.TagRepo,
	commentRepo repositories.CommentRepo,
) *ProjectRepo {
	return &ProjectRepo{
		PG:          PG,
		TagRepo:     tagRepo,
		CommentRepo: commentRepo,
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

	query := "SELECT id,owner_id,title,description,label,selected_bid_id,status,duration,start_time,end_time FROM project WHERE id = $1"

	var duration time.Time
	var start, end sql.NullTime
	var selectedBid sql.NullInt32
	err := repo.PG.QueryRow(ctx, query, projectID).Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &project.Label, &selectedBid, &project.State, &duration, &start, &end)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if selectedBid.Valid {
		project.SelectedBid = int(selectedBid.Int32)
	} else {
		project.SelectedBid = 0
	}

	if start.Valid {
		project.StartTime = start.Time
	}

	if end.Valid {
		project.EndTime = end.Time
	}

	project.Duration = duration

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	tags := repo.TagRepo.GetProjectTag(projectID)

	project.Tags = tags

	comment, _ := repo.CommentRepo.GetCommentOfProject(projectID)

	project.Comment = comment

	return &project, nil
}

func (repo *ProjectRepo) GetUserProject(userID, offset, limit int) []models.ProjectModel {
	var projects []models.ProjectModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,owner_id,title,description,label,selected_bid_id,status,duration,start_time,end_time FROM project WHERE owner_id = $1 ORDER BY id OFFSET $2 LIMIT $3"

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
		var start, end sql.NullTime
		var selectedBid sql.NullInt32
		if err := result.Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &project.Label, &selectedBid, &project.State, &duration, &start, &end); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}

		if selectedBid.Valid {
			project.SelectedBid = int(selectedBid.Int32)
		} else {
			project.SelectedBid = 0
		}

		if start.Valid {
			project.StartTime = start.Time
		}

		if end.Valid {
			project.EndTime = end.Time
		}

		project.Duration = duration
		tag := repo.TagRepo.GetProjectTag(project.ID)
		project.Tags = tag
		comment, _ := repo.CommentRepo.GetCommentOfProject(project.ID)
		project.Comment = comment
		projects = append(projects, project)
	}

	return projects
}

func (repo *ProjectRepo) GetAllProjectsRelatedToUser(userID int) []models.ProjectModel {
	var projects []models.ProjectModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT p.id,p.owner_id,p.title,p.description,p.label,p.selected_bid_id,p.status,p.duration,p.created_time ,p.end_time 
	FROM users AS u 
	JOIN users_team AS ut ON u.id = ut.user_id
	JOIN team AS t ON ut.team_id = t.id
	JOIN bid AS b ON t.id = b.team_id
	JOIN project AS p ON b.project_id = p.id
	WHERE u.id = $1 AND p.status > 2 AND (ut.joined_at < b.created_time) AND (ut.left_at IS NULL OR (p.status = 4 AND ut.left_at < p.end_time) ));`

	result, err := repo.PG.Query(ctx, query, userID)

	if err != nil {
		log.Println("ProjectError: fetching participated projects for user", userID, "details:", err)
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
		var start, end sql.NullTime
		var selectedBid sql.NullInt32
		if err := result.Scan(&project.ID, &project.OwnerID, &project.Title, &project.Description, &project.Label, &selectedBid, &project.State, &duration, &start, &end); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}

		if selectedBid.Valid {
			project.SelectedBid = int(selectedBid.Int32)
		} else {
			project.SelectedBid = 0
		}

		if start.Valid {
			project.StartTime = start.Time
		}

		if end.Valid {
			project.EndTime = end.Time
		}

		project.Duration = duration
		tag := repo.TagRepo.GetProjectTag(project.ID)
		project.Tags = tag
		comment, _ := repo.CommentRepo.GetCommentOfProject(project.ID)
		project.Comment = comment
		projects = append(projects, project)
	}

	return projects
}

func (repo *ProjectRepo) GetProjectCount(userID int) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	query := "SELECT COUNT(*) FROM project WHERE owner_id = $1"

	err := repo.PG.QueryRow(ctx, query, userID).Scan(&count)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return count
}

func (repo *ProjectRepo) CreateProject(ctx context.Context, tx transaction.Tx, userID, label int, title, description string, duration time.Time) int {
	var project_id int

	now := time.Now()

	query := "INSERT INTO project (owner_id, title, description, label, status, duration, created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	row := tx.QueryRow(ctx, query, userID, title, description, label, 1, duration, now, now).(pgx.Row)
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

	now := time.Now()

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

func (repo *ProjectRepo) UpdateProjectState(projectID int, status int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE project SET status=$1 WHERE id=$2"

	_, err := repo.PG.Exec(ctx, query, status, projectID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (repo *ProjectRepo) EndProject(ctx context.Context, tx transaction.Tx, projectID int) {
	query := "UPDATE project SET end_time = CURRENT_TIMESTAMP WHERE id = $1"

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
