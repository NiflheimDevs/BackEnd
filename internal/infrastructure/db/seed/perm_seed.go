package seed

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

/*
ALTER SEQUENCE permission_id_seq RESTART WITH 0;

INSERT INTO "permission" ("name" , "description") VALUES
("ADD_MEMBER", "adds member"),
("REMOVE_MEMEBER", "removes a member"),
("EDIT_INFO", "edit title, bio and ..."),
("BIDDER", "the one who bids"),
("EDIT_NICKNAME", "for teams, it works for positions. for groups and etc for nickname"),
("EDIT_ROLE", "able to change the roles");
*/

func SeedPerms(tx transaction.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()

	query := `SELECT id, name FROM permission`
	rows, err := tx.Query(ctx, query)

	if err != nil {
		panic(err)
	}
	result := rows.(pgx.Rows)

	previousPerms := make(map[string]uint)
	for result.Next() {
		var temp models.PermissionModel
		err = result.Scan(&temp.ID, &temp.Name)
		if err != nil {
			panic(err)
		}
		previousPerms[temp.Name] = temp.ID
		if !enums.PermExists(temp.Name) {
			query = `DELETE FROM role_permission
			WHERE permnission_id = $1`
			tx.Exec(ctx, query, temp.ID)
		}
	}

	query = `DELETE FROM permission`
	_, err = tx.Exec(ctx, query)

	if err != nil {
		panic(err)
	}

	for _, seed := range enums.GetAllPerms() {

		query = `INSERT INTO permission (id,name) VALUES ($1,$2)`
		_, err := tx.Exec(ctx, query, seed, seed.String())

		if err != nil {
			panic(err)
		}
		if previousPerms[seed.String()] != 0 {
			query = `UPDATE role_permission
				SET permission_id = $1
				WHERE permission_id = $2`
			_, err := tx.Exec(ctx, query, seed, previousPerms[seed.String()])

			if err != nil {
				panic(err)
			}
		}
	}

}
