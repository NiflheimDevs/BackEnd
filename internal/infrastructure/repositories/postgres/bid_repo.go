package repositoriesimpl

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type BidRepo struct {
	PG *pgxpool.Pool
}

func NewBidRepo(pg *pgxpool.Pool) *BidRepo {
	return &BidRepo{
		PG: pg,
	}
}

func (br *BidRepo) GetBidInfo(bidID int) (*models.BidModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bid models.BidModel

	query := "SELECT id,team_id,project_id,prepayment,total,description,expected_time,created_time FROM bid WHERE id = $1"

	err := br.PG.QueryRow(ctx, query, bidID).Scan(&bid.ID, &bid.TeamID, &bid.ProjectID, &bid.PrePayment, &bid.Total, &bid.Description, &bid.ExpectedTime, &bid.CreatedTime)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return &bid, nil
}

func (br *BidRepo) GetTeamBids(teamID int64) []models.BidModel {
	var bids []models.BidModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,team_id,project_id,prepayment,total,description,expected_time,created_time FROM bid WHERE team_id = $1"

	results, err := br.PG.Query(ctx, query, teamID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
	defer results.Close()
	for results.Next() {
		var bid models.BidModel
		var time2 time.Time
		if err := results.Scan(&bid.ID, &bid.TeamID, &bid.ProjectID, &bid.PrePayment, &bid.Total, &bid.Description, &bid.ExpectedTime, &time2); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		bid.CreatedTime = time2
		bids = append(bids, bid)
	}
	return bids
}

func (br *BidRepo) GetBidOfProject(projectID int) []models.BidModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,team_id,project_id,prepayment,total,description,expected_time,created_time FROM bid WHERE project_id = $1"

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
		var time time.Time
		var des sql.NullString
		if err := results.Scan(&bid.ID, &bid.TeamID, &bid.ProjectID, &bid.PrePayment, &bid.Total, &des, &bid.ExpectedTime, &time); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		bid.CreatedTime = time
		if des.Valid {
			bid.Description = des.String
		}

		bids = append(bids, bid)
	}
	return bids
}

func (br *BidRepo) PutBid(info dto.BidInfo) int {
	var bidid int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO bid(team_id,project_id,prepayment,total,description,expected_time) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"

	err := br.PG.QueryRow(ctx, query, info.TeamID, info.ProjectID, info.PP, info.Total, info.Description, info.ExpectedTime).Scan(&bidid)

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

func (br *BidRepo) AcceptBid(ctx context.Context, tx transaction.Tx, bidID int, projectID int) {
	query := "UPDATE project SET selected_bid_id = $1,status = 3 WHERE id = $2"

	_, err := br.PG.Exec(ctx, query, bidID, projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (br *BidRepo) UpdateBid(info dto.BidInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE bid SET team_id = $1, project_id = $2, prepayment = $3, total = $4, description = $5, expected_time = $6 WHERE id = $7`

	_, err := br.PG.Exec(ctx, query, info.TeamID, info.ProjectID, info.PP, info.Total, info.Description, info.ExpectedTime, info.BidID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}
