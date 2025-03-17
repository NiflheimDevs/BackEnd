package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (repo *ProjectRepo) CreateProject(userID int, title, description, duration string) (int, error) {
	var project_id int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO project (owner_id, title, description, duration) VALUES ($1, $2, $3, $4) RETURNING id"

	err := repo.PG.QueryRow(ctx, query, userID, title, description, duration).Scan(&project_id)

	return project_id, err
}

func (repo *ProjectRepo) AddProjectTag(projectID, tagID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users_project_tag (project_user_id, tag_id, type) VALUES ($1, $2, 2)"

	_, err := repo.PG.Exec(ctx, query, projectID, tagID)

	return err
}
