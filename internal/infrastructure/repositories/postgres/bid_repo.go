package repositoriesimpl

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidRepo struct {
	PG *pgxpool.Pool
}

func NewBidRepo(pg *pgxpool.Pool) *BidRepo {
	return &BidRepo{
		PG: pg,
	}
}

func (br *BidRepo) GetBidOfProject(projectID int) []models.BidModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,team_id,project_id,value,expected_time,created_time FROM bid WHERE project_id = $1"

	results, err := br.PG.Query(ctx, query, projectID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
	var bids []models.BidModel
	defer results.Close()
	for results.Next() {
		var bid models.BidModel
		var time1 time.Time
		var time2 time.Time
		if err := results.Scan(&bid.ID, &bid.TeamID, &bid.ProjectID, &bid.Value, &time1, &time2); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		bid.ExpectedTime = time1.Format("2006-01-02 15:04:05")
		bid.CreatedTime = time2.Format("2006-01-02 15:04:05")
		bids = append(bids, bid)
	}
	return bids
}

func (br *BidRepo) PutBid(teamID int, projectID int, value int64, expected_time time.Time) int {
	var bidid int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO bid(team_id,project_id,value,expected_time) VALUE ($1,$2,$3,$4) RETURNING id"

	err := br.PG.QueryRow(ctx, query, teamID, projectID, value, expected_time.Format("2006-01-02 15:04:05")).Scan(&bidid)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return bidid
}
