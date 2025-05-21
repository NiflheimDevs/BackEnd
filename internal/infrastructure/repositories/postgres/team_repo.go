package repositoriesimpl

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
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

	row := tx.QueryRow(ctx, query, strconv.Itoa(userid)).(pgx.Row)
	err := row.Scan(&teamid)

	if err != nil {
		log.Println("TeamError: user", userid, "can't have a team! error detail:", err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
	return teamid
}

func (tr *TeamRepo) AddMember(userid int, teamid int64, position string, roleid enums.RoleType) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `INSERT INTO users_team
    (user_id , team_id, position, role_id) VALUES
    ($1, $2, $3, $4)`

	_, err := tr.PG.Exec(ctx, query, userid, teamid, position, roleid)

	return err
}

func (tr *TeamRepo) AddMemberWithTx(ctx context.Context, tx transaction.Tx, userid int, teamid int64, position string, roleid enums.RoleType) error {

	query := `INSERT INTO users_team
    (user_id , team_id, position, role_id) VALUES
    ($1, $2, $3, $4)`

	_, err := tx.Exec(ctx, query, userid, teamid, position, roleid)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
			})
		}
	}

	return err
}

func (tr *TeamRepo) RemoveMember(userid int, teamid int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `UPDATE users_team 
	SET left_at = CURRENT_TIMESTAMP 
	WHERE user_id = $1 AND team_id = $2`

	_, err := tr.PG.Exec(ctx, query, userid, teamid)

	return err
}

func (tr *TeamRepo) UpdateMemberPosition(info *dto.UpdateMemberPositionDto) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `UPDATE users_team
	SET position = $1 
	WHERE user_id = $2 AND team_id = $3 AND left_at IS NULL`

	_, err := tr.PG.Exec(ctx, query, info.NewPosition, info.Userid, info.Teamid)

	return err
}

// TODO: handle this shit
func (tr *TeamRepo) UpdateMemberRole(newRole enums.RoleType, userid int, teamid int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `UPDATE users_team
	SET role_id = $1 
	WHERE user_id = $2 AND team_id = $3 AND left_at IS NULL`

	_, err := tr.PG.Exec(ctx, query, newRole, userid, teamid)

	return err
}

func (tr *TeamRepo) CreateTeam(ctx context.Context, tx transaction.Tx, title string, description string) int64 {

	query := `INSERT INTO team
    (title , description) VALUES
    ($1, $2) 
    RETURNING id`

	var teamid int64

	row := tx.QueryRow(ctx, query, title, description).(pgx.Row)
	err := row.Scan(&teamid)

	if err != nil {
		log.Println("TeamError: error detail:", err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
	return teamid
}

func (tr *TeamRepo) UpdateTeam(teamid int64, title string, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	query := `UPDATE team
	SET title = $1, description = $2
	WHERE type = 0 AND id = $3`

	_, err := tr.PG.Exec(ctx, query, title, description, teamid)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
			})
		}
	}

	return err
}

func (tr *TeamRepo) DeleteTeam(teamid int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `DELETE FROM team
	WHERE id = $1 AND type = 0`

	_, err := tr.PG.Exec(ctx, query, teamid)

	return err
}

func (tr *TeamRepo) GetEveryTeamID(userID int) []int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT t.id FROM team AS t JOIN users_team AS ut ON t.id = ut.team_id WHERE ut.user_id = $1"

	result, err := tr.PG.Query(ctx, query, userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	var ids []int64
	defer result.Close()
	for result.Next() {
		var id int64
		err := result.Scan(&id)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		ids = append(ids, id)
	}

	return ids
}

func (tr *TeamRepo) GetEveryTeamInfo(teamid int64) *models.TeamModel {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT t.id, t.title, t.type
		FROM team as t
		WHERE t.id = $1`

	var team models.TeamModel

	err := tr.PG.QueryRow(ctx, query, teamid).Scan(&team.ID, &team.Title, &team.Type)

	if err != nil {
		log.Println("TeamError: Team id", teamid, "not found. details:", err)
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil
	}
	return &team
}

func (tr *TeamRepo) GetTeam(teamid int64) *models.TeamModel {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT t.id, t.title, t.description, t.created_at
		FROM team as t
		WHERE t.id = $1 AND t.type = 0`

	var team models.TeamModel
	var description sql.NullString

	err := tr.PG.QueryRow(ctx, query, teamid).Scan(&team.ID, &team.Title, &description, &team.Created_at)

	if err != nil {
		log.Println("TeamError: Team id", teamid, "not found. details:", err)
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil
	}
	if description.Valid {
		team.Description = description.String
	}
	return &team

}

func (tr *TeamRepo) GetTeamForUser(userid int, teamid int64) *models.TeamModel {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT t.id, t.title, t.description, t.created_at
		FROM team as t
		JOIN users_team as ut
		WHERE ut.user_id = $1 AND ut.team_id = $2 AND t.type = 0 AND ut.left_at IS NULL`

	var team models.TeamModel
	var description sql.NullString

	err := tr.PG.QueryRow(ctx, query, userid, teamid).Scan(&team.ID, &team.Title, &description, &team.Created_at)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		return nil
	}
	if description.Valid {
		team.Description = description.String
	}
	return &team

}

func (tr *TeamRepo) GetTeamsForUser(userid int, isActive, dontCare bool) []dto.GetTeamPreviewDto {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	var query string
	if dontCare {
		query = `SELECT t.id, t.title, t.description, ut.position, ut.joined_at , ut.left_at
		FROM team as t
		JOIN users_team as ut
		ON t.id = ut.team_id
		WHERE ut.user_id = $1 AND t.type = 0`
	} else {
		if isActive {
			query = `SELECT t.id, t.title, t.description, ut.position, ut.joined_at , ut.left_at
		FROM team as t
		JOIN users_team as ut
		ON t.id = ut.team_id
		WHERE ut.user_id = $1 AND t.type = 0 AND ut.left_at IS NULL`
		} else {
			query = `SELECT t.id, t.title, t.description, ut.position, ut.joined_at , ut.left_at
		FROM team as t
		JOIN users_team as ut
		ON t.id = ut.team_id
		WHERE ut.user_id = $1 AND t.type = 0 AND ut.left_at IS NULL`
		}
	}
	rows, err := tr.PG.Query(ctx, query, userid)
	if err != nil {
		log.Println("TeamError: error fetching teams for user", userid, "error detail:", err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
	defer rows.Close()

	var teams []dto.GetTeamPreviewDto
	for rows.Next() {
		var team dto.GetTeamPreviewDto
		var description, position sql.NullString
		var leftAt sql.NullTime

		if err := rows.Scan(&team.ID, &team.Title, &description, &position, &team.JoinedAt, &leftAt); err != nil {
			log.Println("TeamError: error scanning team for user", userid, "error detail:", err)
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		if description.Valid {
			team.Description = description.String
		}
		if position.Valid {
			team.Position = position.String
		}
		if leftAt.Valid {
			team.LeftAt = leftAt.Time
		}
		teams = append(teams, team)
	}

	return teams
}

func (tr *TeamRepo) GetTeamsForUserWithRole(userid int) []dto.GetTeamWithRole {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT t.id, t.title, t.description, ut.position, ut.role_id
		FROM team as t
		JOIN users_team as ut
		ON t.id = ut.team_id
		WHERE ut.user_id = $1 AND t.type = 0 AND ut.left_at IS NULL`

	rows, err := tr.PG.Query(ctx, query, userid)
	if err != nil {
		log.Println("TeamError: error fetching teams for user", userid, "error detail:", err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
	defer rows.Close()

	var teams []dto.GetTeamWithRole
	for rows.Next() {
		var team dto.GetTeamWithRole
		var description, position sql.NullString

		if err := rows.Scan(&team.ID, &team.Title, &description, &position, &team.RoleId); err != nil {
			log.Println("TeamError: error scanning team for user", userid, "error detail:", err)
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		if description.Valid {
			team.Description = description.String
		}
		if position.Valid {
			team.Position = position.String
		}
		teams = append(teams, team)
	}

	return teams
}

func (tr *TeamRepo) GetTeamInfo(teamid int64) (*dto.GetInternalTeamInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	var team dto.GetInternalTeamInfo
	var description sql.NullString

	query := `SELECT t.id, t.title, t.description
		FROM team as t
		WHERE t.id = $1 AND t.type = 0`

	err := tr.PG.QueryRow(ctx, query, teamid).Scan(&team.ID, &team.Title, &description)
	if description.Valid {
		team.Description = description.String
	}
	if err == pgx.ErrNoRows {
		return nil, err
	}
	if err != nil {
		log.Println("TeamError: error fetching teams info", teamid, "error detail:", err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	owner := tr.GetMembersForTeamFilterdByRole(teamid, enums.TEAM_OWNER)

	if len(owner) > 0 {
		team.OwnerID = owner[0].Info.Userid
	}

	return &team, nil
}

func (tr *TeamRepo) GetOneManTeamInfo(teamid int64) (*dto.GetInternalTeamInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	firstquery := `
		SELECT t.title
		FROM team AS t
		WHERE t.id = $1 AND t.type = 1
	`

	secondquery := `SELECT u.id, u.username, u.bio
					FROM users AS u
					WHERE u.id = $1`

	var team dto.GetInternalTeamInfo
	var useridstring string

	err := tr.PG.QueryRow(ctx, firstquery, teamid).Scan(&useridstring)
	if err == pgx.ErrNoRows {
		return nil, err
	}

	userid, _ := strconv.Atoi(useridstring)

	var bio sql.NullString
	err = tr.PG.QueryRow(ctx, secondquery, userid).Scan(&team.OwnerID, &team.Title, &bio)
	if err == pgx.ErrNoRows {
		return nil, err
	}
	if err != nil {
		log.Println("TeamError: error fetching user info", teamid, "error detail:", err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	if bio.Valid {
		team.Description = bio.String
	}
	team.ID = int64(userid)

	return &team, nil
}

func (tr *TeamRepo) GetMembersForTeam(teamid int64) []dto.ReadMemberDto {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT u.id, u.username, u.firstname, u.lastname, ut.position, ut.role_id
		FROM users_team as ut
		JOIN users AS u
		ON ut.user_id = u.id 
		WHERE ut.team_id = $1 AND ut.left_at IS NULL`

	rows, err := tr.PG.Query(ctx, query, teamid)
	if err != nil {
		log.Println("TeamError: error fetching members for team", teamid, "error detail:", err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
	defer rows.Close()

	var members []dto.ReadMemberDto

	for rows.Next() {
		var member dto.ReadMemberDto
		var info dto.MemberInfoDto
		var role sql.NullInt64
		var firstname, lastname, position sql.NullString
		if err := rows.Scan(&info.Userid, &info.Username, &firstname, &lastname, &position, &role); err != nil {
			log.Println("TeamError: error scanning users for team", teamid, "error detail:", err)
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}

		member.Info = &info
		if role.Valid {
			member.Role = enums.RoleType(role.Int64)
		}
		if firstname.Valid {
			info.FirstName = firstname.String
		}
		if lastname.Valid {
			info.LastName = lastname.String
		}
		if position.Valid {
			info.Position = position.String
		}

		members = append(members, member)
	}

	return members
}

func (tr *TeamRepo) GetMembersForTeamFilterdByRole(teamid int64, roleid enums.RoleType) []dto.ReadMemberDto {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT u.id, u.username, u.firstname, u.lastname, ut.position
		FROM users_team as ut
		JOIN users AS u
		ON ut.user_id = u.id 
		WHERE ut.team_id = $1 AND ut.role_id = $2 AND ut.left_at IS NULL`

	rows, err := tr.PG.Query(ctx, query, teamid, roleid)
	if err != nil {
		log.Println("TeamError: error fetching members for team", teamid, "error detail:", err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
	defer rows.Close()

	var members []dto.ReadMemberDto

	for rows.Next() {
		var member dto.ReadMemberDto
		var info dto.MemberInfoDto
		var firstname, lastname, position sql.NullString
		if err := rows.Scan(&info.Userid, &info.Username, &firstname, &lastname, &position); err != nil {
			log.Println("TeamError: error scanning users for team", teamid, "error detail:", err)
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}

		member.Info = &info
		member.Role = roleid
		if firstname.Valid {
			info.FirstName = firstname.String
		}
		if lastname.Valid {
			info.LastName = lastname.String
		}
		if position.Valid {
			info.Position = position.String
		}

		members = append(members, member)
	}

	return members
}

func (tr *TeamRepo) GetMemberForTeam(teamid int64, userid int) *dto.ReadMemberDto {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	query := `SELECT u.id, u.username, u.firstname, u.lastname, ut.position, ut.role_id
		FROM users_team as ut
		JOIN users AS u
		ON ut.user_id = u.id 
		WHERE ut.team_id = $1 AND u.id = $2 AND ut.left_at IS NULL`
	var member dto.ReadMemberDto
	var info dto.MemberInfoDto
	var role sql.NullInt64
	var firstname, lastname, position sql.NullString

	if err := tr.PG.QueryRow(ctx, query, teamid, userid).Scan(&info.Userid, &info.Username, &firstname, &lastname, &position, &role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println(err, teamid, userid)
			return nil
		} else {
			log.Println("TeamError: error scanning users for team", teamid, "error detail:", err)
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
	}

	member.Info = &info
	if role.Valid {
		member.Role = enums.RoleType(role.Int64)
	}
	if firstname.Valid {
		info.FirstName = firstname.String
	}
	if lastname.Valid {
		info.LastName = lastname.String
	}
	if position.Valid {
		info.Position = position.String
	}
	return &member
}

func (tr *TeamRepo) GetOneManTeamID(userid int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var teamid int64

	query := "SELECT id FROM team WHERE title = $1"

	err := tr.PG.QueryRow(ctx, query, strconv.Itoa(userid)).Scan(&teamid)
	if err == pgx.ErrNoRows {
		return 0, err
	}
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
	return teamid, nil
}
