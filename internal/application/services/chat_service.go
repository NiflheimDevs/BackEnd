package services

import "github.com/niflheimdevs/backend/internal/application/dto"

type ChatService interface {
	GetAllRoom(userID int) []dto.RoomInfo
	GetRoomMessages(userID, roomID int) []dto.Message
	CreateUserRoom(userID int, targetUserID int)
	SaveMessage(roomID int, senderID int, content string) *dto.Message
}
