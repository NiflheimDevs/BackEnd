package seed

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

func SeedRoles(tx transaction.Tx) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()

	query := `SELECT id, name FROM role`
	rows, err := tx.Query(ctx, query)

	if err != nil {
		panic(err)
	}
	result := rows.(pgx.Rows)

	previousRoles := make(map[string]uint64)
	for result.Next() {
		var temp models.RoleModel
		err = result.Scan(&temp.ID, &temp.Name)
		if err != nil {
			panic(err)
		}
		previousRoles[temp.Name] = temp.ID
		if !enums.RoleExists(temp.Name) {
			query = `DELETE FROM role_permission
			WHERE role_id = $1`
			_, err = tx.Exec(ctx, query, temp.ID)
			if err != nil {
				panic(err)
			}
			query = `DELETE FROM users_team
			WHERE role_id = $1`
			_, err = tx.Exec(ctx, query, temp.ID)
			if err != nil {
				panic(err)
			}
			query = `DELETE FROM users_chat
			WHERE role_id = $1`
			_, err = tx.Exec(ctx, query, temp.ID)
			if err != nil {
				panic(err)
			}
		}
	}

	query = `DELETE FROM role`
	_, err = tx.Exec(ctx, query)

	if err != nil {
		panic(err)
	}

	query = `ALTER SEQUENCE role_id_seq RESTART WITH 0`
	_, err = tx.Exec(ctx, query)

	if err != nil {
		panic(err)
	}
	for _, seed := range enums.GetAllRoles() {

		query = `INSERT INTO role (id,name) VALUES ($1,$2)`
		_, err = tx.Exec(ctx, query, seed, seed.String())
		if err != nil {
			panic(err)
		}
		if previousRoles[seed.String()] != 0 {

			query = `UPDATE uesrs_team
				SET role_id = $1
				WHERE role_id = $2`
			_, err = tx.Exec(ctx, query, seed, previousRoles[seed.String()])
			if err != nil {
				panic(err)
			}

			query = `UPDATE users_chat
				SET role_id = $1
				WHERE role_id = $2`
			_, err = tx.Exec(ctx, query, seed, previousRoles[seed.String()])
			if err != nil {
				panic(err)
			}
		}
		// permission for roles
		query = `INSERT INTO role_permission 
			(role_id, permission_id) VALUES
			($1 , $2)`
		for _, perm := range seed.GetPermissionsForRole() {
			_, err = tx.Exec(ctx, query, seed, perm)
			if err != nil {
				panic(err)
			}
		}
	}

}
