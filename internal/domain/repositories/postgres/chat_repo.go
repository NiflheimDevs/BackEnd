package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type ChatRepo interface {
	CreateChatRoom(ctx context.Context, tx transaction.Tx, title, description string) int
	InsertMemberToChat(ctx context.Context, tx transaction.Tx, userID, chatID int)
	GetAllRoomInfo(userID int) []dto.RoomInfo
	GetRoomInfo(roomID int) (*models.ChatModel, error)
	GetRoomMessages(roomID int) []models.MessageModel
	SaveMessage(roomID int, userID int, content string) int
}
