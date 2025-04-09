package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type LabelRepo struct {
	PG *pgxpool.Pool
}

func NewLabelRepo(pg *pgxpool.Pool) *LabelRepo {
	return &LabelRepo{
		PG: pg,
	}
}
func (lr *LabelRepo) GetLabels() []models.LabelModel {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM label`
	result, err := lr.PG.Query(ctx, query)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	var labels []models.LabelModel
	for result.Next() {
		var label models.LabelModel
		err = result.Scan(&label.ID, &label.Name, &label.Description, &label.Price)
		if err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		labels = append(labels, label)
	}
	return labels
}

func (lr *LabelRepo) GetLabelInfo(labelID int) *models.LabelModel {
	var label models.LabelModel

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id,name,description,price FROM label WHERE id = $1`
	err := lr.PG.QueryRow(ctx, query, labelID).Scan(&label.ID, &label.Name, &label.Description, &label.Price)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return &label
}
