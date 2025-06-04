package repositoriesimpl

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type NotifRepo struct {
	PG *pgxpool.Pool
}

func NewNotifRepo(pg *pgxpool.Pool) *NotifRepo {
	return &NotifRepo{
		PG: pg,
	}
}

func (nr *NotifRepo) SaveNotif(userID int, message string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int

	query := "INSERT INTO notification(user_id,content) VALUES ($1,$2)"

	err := nr.PG.QueryRow(ctx, query, userID, message).Scan(&id)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return id
}

func (nr *NotifRepo) GetNotifs(userID int) []models.NotifModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select id,user_id,content,send_time,read_time,isread from notification where user_id = $1"

	result, err := nr.PG.Query(ctx, query, userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	var notifs []models.NotifModel

	defer result.Close()
	for result.Next() {
		var notif models.NotifModel
		var read sql.NullTime
		err := result.Scan(&notif.ID, &notif.UserID, &notif.Content, &notif.SendTime, &read, &notif.IsRead)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		if read.Valid {
			notif.ReadTime = read.Time
		}
		notifs = append(notifs, notif)
	}

	return notifs
}

func (nr *NotifRepo) GetNotif(id int) models.NotifModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select id,user_id,content,send_time,read_time,isread from notification where id = $1"

	var notif models.NotifModel
	var read sql.NullTime

	err := nr.PG.QueryRow(ctx, query, id).Scan(&notif.ID, &notif.UserID, &notif.Content, &notif.SendTime, &read, &notif.IsRead)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	if read.Valid {
		notif.ReadTime = read.Time
	}
	return notif
}

func (nr *NotifRepo) ReadNotif(id int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	now := time.Now()

	query := "update notification set read_time = $1 , isread = true where id = $2"

	_, err := nr.PG.Exec(ctx, query, now, id)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
}
