package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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
