package repositoriesimpl

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TeamRepo struct {
	PG *pgxpool.Pool
}

func NewTeamRepo(
	PG *pgxpool.Pool,
) *TeamRepo {
	return &TeamRepo{
		PG: PG,
	}
}

// takes userid only to make the title of the team point to the user. ( avoiding a join if everything goes right )
func (tr *TeamRepo) CreateOneManTeam(ctx context.Context, tx transaction.Tx, userid int) int64 {

	query := `INSERT INTO team
    (type, title) VALUES
    (1, $1) 
    RETURNING id`

	var teamid int64

	row := tx.QueryRow(ctx, query, userid).(pgx.Row)
	err := row.Scan(&teamid)

	if err != nil {
		log.Println("TeamError: user", userid, "can't have a team! error detail:", err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
	return teamid
}

func (tr *TeamRepo) AddMember(userid int, teamid int64, position string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `INSERT INTO users_team
    (user_id , team_id, position) VALUES
    ($1, $2, $3)`

	_, err := tr.PG.Exec(ctx, query, userid, teamid, position)

	return err
}

func (tr *TeamRepo) AddMemberWithTx(ctx context.Context, tx transaction.Tx, userid int, teamid int64, position string) error {

	query := `INSERT INTO users_team
    (user_id , team_id, position) VALUES
    ($1, $2, $3)`

	_, err := tx.Exec(ctx, query, userid, teamid, position)

	return err
}
