package repositoriesimpl

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type RoleRepo struct {
	PG *pgxpool.Pool
}

func NewRoleRepo(
	PG *pgxpool.Pool,
) *RoleRepo {
	return &RoleRepo{
		PG: PG,
	}
}

// type_ : 0 for user, 1 for team and 2 for chat
func (rr *RoleRepo) AddRole(userid int, originid int, type_ int, roleid uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `INSERT INTO member_role
    (user_id , origin_id, type, role_id) VALUES
    ($1, $2, $3, $4)`

	_, err := rr.PG.Exec(ctx, query, userid, originid, type_, roleid)

	return err
}

// type_ : 0 for user, 1 for team and 2 for chat
func (rr *RoleRepo) AddRoleWithTx(ctx context.Context, tx transaction.Tx, userid int, originid int, type_ int, roleid uint) error {

	query := `INSERT INTO member_role
    (user_id , origin_id, type, role_id) VALUES
    ($1, $2, $3, $4)`

	_, err := tx.Exec(ctx, query, userid, originid, type_, roleid)

	return err
}

func (rr *RoleRepo) UpdateRole(userid int, originid int, type_ int, roleid uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `UPDATE member_role
    SET role_id = $1 
    WHERE 
    user_id = $2 AND type = $3 AND origin_id = $4`

	_, err := rr.PG.Exec(ctx, query, roleid, userid, type_, originid)

	return err
}

func (rr *RoleRepo) DeleteRole(userid int, originid int, type_ int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `DELETE FROM member_role
    WHERE 
    user_id = $1 AND type = $3 AND origin_id = $2`

	_, err := rr.PG.Exec(ctx, query, userid, originid, type_)

	return err
}

func (rr *RoleRepo) GetRoleModel(userid int, originid int, type_ int) *models.RoleModel {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT r.id , r.name 
    FROM member_role AS mr
    JOIN role AS r 
    ON r.id = mr.role_id
    WHERE mr.user_id = $1 AND mr.origin_id = $2 AND mr.type = $3`

	var theRole models.RoleModel

	err := rr.PG.QueryRow(ctx, query, userid, originid, type_).Scan(&theRole.ID, &theRole.Name)

	if err != nil {
		return nil
	}
	return &theRole
}

func (rr *RoleRepo) GetRoleId(name string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT r.id
    FROM role AS r 
    WHERE r.name = $1`

	var id uint

	err := rr.PG.QueryRow(ctx, query, name).Scan(&id)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return 0, err
	}
	return id, nil
}

func (rr *RoleRepo) GetRoleName(id uint) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT r.name
    FROM role AS r 
    WHERE r.id = $1`

	var name string

	err := rr.PG.QueryRow(ctx, query, name).Scan(&name)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return "", err
	}
	return name, nil
}

func (rr *RoleRepo) GetPermissionsForRole(roleid uint) []models.PermissionModel {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	query := `SELECT p.id, p.name 
	FROM role_permission AS rp
	JOIN permission AS p 
	ON p.id = rp.permission_id
	WHERE rp.role_id = $1`
	result, err := rr.PG.Query(ctx, query)

	if err != nil {
		log.Println("PermissionError: ", err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

	var perms []models.PermissionModel
	for result.Next() {
		var perm models.PermissionModel
		err = result.Scan(&perm.ID, &perm.Name)
		if err != nil {
			log.Println("PermissionError on parsing: ", err)
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
			})
		}
		perms = append(perms, perm)
	}

	return perms
}
