package repositoriesimpl

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type CommentRepo struct {
	PG *pgxpool.Pool
}

func NewCommentRepo(pg *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{
		PG: pg,
	}
}

func (cr *CommentRepo) AddComment(projectID, bidID int, content string, star int) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO comment(project_id,bid_id,content,rating) VALUES ($1,$2,$3,$4) RETURNING ID"

	var id int

	err := cr.PG.QueryRow(ctx, query, projectID, bidID, content, star).Scan(&id)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return id
}

func (cr *CommentRepo) GetStar(userID int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select Avg(c.rating)
			  from comment as c
			  join bid as b
			  on c.bid_id = b.id
			  join team as t
			  on b.team_id = t.id
			  join users_team as ut
			  on t.id = ut.team_id
			  join project as p
			  on b.project_id = p.id
			  where (t.type = 1 and t.title = $1) or (ut.user_id = $1::integer and (ut.left_at is null or ut.left_at > p.end_time))`

	var average sql.NullFloat64

	err := cr.PG.QueryRow(ctx, query, strconv.Itoa(userID)).Scan(&average)

	if err == pgx.ErrNoRows {
		return 0, err
	}

	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	if average.Valid {
		return average.Float64, nil
	}

	return 0, nil
}

func (cr *CommentRepo) GetUserComments(userID int) ([]dto.CommentWithUserDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT c.id AS comment_id, c.project_id, c.content, c.rating,
				u.id AS user_id, u.firstname, u.lastname, u.username
			FROM comment AS c
			JOIN bid AS b ON c.bid_id = b.id
			JOIN team AS t ON b.team_id = t.id
			JOIN project AS p ON c.project_id = p.id
			JOIN users AS u ON p.owner_id = u.id
			WHERE t.type = 1 AND t.title = $1

			UNION

			SELECT c.id AS comment_id, c.project_id, c.content, c.rating,
				u.id AS user_id, u.firstname, u.lastname, u.username
			FROM comment AS c
			JOIN bid AS b ON c.bid_id = b.id
			JOIN team AS t ON b.team_id = t.id
			JOIN users_team AS ut ON t.id = ut.team_id
			JOIN project AS p ON c.project_id = p.id
			JOIN users AS u ON p.owner_id = u.id
			WHERE ut.user_id = $2 AND (ut.left_at > p.end_time OR ut.left_at IS NULL)
			`

	results, err := cr.PG.Query(ctx, query, strconv.Itoa(userID), userID)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	var comments []dto.CommentWithUserDTO

	defer results.Close()
	for results.Next() {
		var comment dto.CommentWithUserDTO
		var first, last sql.NullString
		err := results.Scan(&comment.ID, &comment.ProjectID, &comment.Content, &comment.Rating, &comment.UserID, &first, &last, &comment.Username)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		if first.Valid {
			comment.Firstname = first.String
		}
		if last.Valid {
			comment.Lastname = last.String
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (cr *CommentRepo) GetCommentInfo(id int) dto.CommentWithUserDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT c.id AS comment_id, c.project_id, c.content, c.rating,
				u.id AS user_id, u.firstname, u.lastname, u.username
			  FROM comment AS c
			  JOIN bid AS b
			  ON c.bid_id = b.id
			  JOIN project AS p
			  on b.project_id = p.id
			  JOIN users AS u
			  on p.owner_id = u.id
			  WHERE c.id = $1`

	var comment dto.CommentWithUserDTO
	var first, last sql.NullString

	err := cr.PG.QueryRow(ctx, query, id).Scan(&comment.ID, &comment.ProjectID, &comment.Content, &comment.Rating, &comment.UserID, &first, &last, &comment.Username)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
	if first.Valid {
		comment.Firstname = first.String
	}
	if last.Valid {
		comment.Lastname = last.String
	}

	return comment
}

func (cr *CommentRepo) GetCommentOfProject(projectID int) (*models.CommentModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT c.id AS comment_id, c.project_id, c.content, c.rating
			  FROM comment AS c
			  WHERE c.project_id = $1`

	var comment models.CommentModel

	err := cr.PG.QueryRow(ctx, query, projectID).Scan(&comment.ID, &comment.ProjectID, &comment.Content, &comment.Rating)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return &comment, nil
}
