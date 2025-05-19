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

type ChatRepo struct {
	PG *pgxpool.Pool
}

func NewChatRepo(pg *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{
		PG: pg,
	}
}

func (cr *ChatRepo) GetRoomInfo(roomID int) (*models.ChatModel, error) {
	var room models.ChatModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM team WHERE id = $1"

	err := cr.PG.QueryRow(ctx, query, roomID).Scan(&room.RoomID, &room.Title, &room.Description)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return &room, nil
}

func (cr *ChatRepo) GetAllRoomInfo(userID int) []dto.RoomInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT uc.chat_id,u.id,u.username,u.firstname,u.lastname
				FROM users u
				JOIN users_chat uc ON u.id = uc.user_id
				WHERE uc.chat_id IN (
    				SELECT chat_id FROM users_chat WHERE user_id = $1
				)
				AND u.id != $1;`
	results, err := cr.PG.Query(ctx, query, userID)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	var users []dto.RoomInfo

	defer results.Close()
	for results.Next() {
		var user dto.RoomInfo
		var first, last sql.NullString
		err := results.Scan(&user.RoomID, &user.UserID, &user.Username, &first, &last)
		if err != nil {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}
		if first.Valid {
			user.FirstName = first.String
		}
		if last.Valid {
			user.LastName = last.String
		}
		users = append(users, user)
	}
	return users
}

func (cr *ChatRepo) CreateChatRoom(ctx context.Context, tx transaction.Tx, title, description string) int {
	query := "INSERT INTO chat(title,description) VALUES ($1,$2) RETURNING id"

	var id int

	result := tx.QueryRow(ctx, query, title, description).(pgx.Row)
	err := result.Scan(&id)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return id
}

func (cr *ChatRepo) InsertMemberToChat(ctx context.Context, tx transaction.Tx, userID, chatID int) {
	query := "INSERT INTO users_chat(chat_id,user_id,role_id) VALUES ($1,$2,1)"

	_, err := tx.Exec(ctx, query, chatID, userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}
}

func (cr *ChatRepo) GetRoomMessages(roomID int) []models.MessageModel {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * from message AS m WHERE chat_id = $1"

	results, err := cr.PG.Query(ctx, query, roomID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	var messages []models.MessageModel

	defer results.Close()
	for results.Next() {
		var message models.MessageModel

		err := results.Scan(&message.MessageID, &message.RoomID, &message.SenderID, &message.SendTime, &message.EditTime, &message.Content)

		if err != nil {
			panic(exceptions.Exception{
				Tag:    exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
			})
		}

		messages = append(messages, message)
	}

	return messages
}

func (cr *ChatRepo) SaveMessage(roomID int, userID int, content string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int

	query := "INSERT INTO message(chat_id,sender_id,content) VALUE ($1,$2,$3) RETURNING id"

	err := cr.PG.QueryRow(ctx, query, roomID, userID, content).Scan(&id)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.DATABASE_ERROR},
		})
	}

	return id
}
