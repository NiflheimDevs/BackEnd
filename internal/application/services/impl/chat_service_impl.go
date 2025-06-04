package servicesimpl

import (
	"context"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type ChatService struct {
	TxManager   transaction.TxManager
	ChatRepo    repositories.ChatRepo
	UserService services.UserService
}

func NewChatService(
	txManager transaction.TxManager,
	chatRepo repositories.ChatRepo,
	userService services.UserService,
) *ChatService {
	return &ChatService{
		TxManager:   txManager,
		ChatRepo:    chatRepo,
		UserService: userService,
	}
}

func (cs *ChatService) GetAllRoom(userID int) []dto.RoomInfo {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
	}

	rooms := cs.ChatRepo.GetAllRoomInfo(userID)

	return rooms
}

func (cs *ChatService) GetRoomMessages(userID, roomID int) []dto.Message {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
	}

	_, err := cs.ChatRepo.GetRoomInfo(roomID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{exceptions.CHAT_NOT_FOUND},
		})
	}

	messages := cs.ChatRepo.GetRoomMessages(roomID)

	var messagesDTO []dto.Message

	for _, message := range messages {
		messageDTO := dto.Message{
			ID:       message.MessageID,
			Content:  message.Content,
			SenderID: message.SenderID,
			SendTime: message.SendTime,
		}
		if message.SenderID == userID {
			messageDTO.Type = 1
		} else {
			messageDTO.Type = 2
		}
		messagesDTO = append(messagesDTO, messageDTO)
	}

	return messagesDTO
}

func (cs *ChatService) CreateUserRoom(userID int, targetUserID int) *dto.RoomInfo {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
	}

	userinfo := cs.UserService.GetUserInfo(targetUserID, userID)

	rooms := cs.GetAllRoom(userID)

	for _, room := range rooms {
		if room.UserID == targetUserID {
			return &room
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := cs.TxManager.Begin(ctx)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	roomID := cs.ChatRepo.CreateChatRoom(ctx, tx, "", "")
	cs.ChatRepo.InsertMemberToChat(ctx, tx, userID, roomID)
	cs.ChatRepo.InsertMemberToChat(ctx, tx, targetUserID, roomID)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	roomInfo := &dto.RoomInfo{
		FirstName: userinfo.FirstName,
		LastName:  userinfo.LastName,
		Username:  userinfo.Username,
		UserID:    targetUserID,
		RoomID:    roomID,
	}

	return roomInfo
}

func (cs *ChatService) SaveMessage(roomID int, senderID int, content string) *dto.Message {
	if senderID == -1 || senderID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
	}

	id := cs.ChatRepo.SaveMessage(roomID, senderID, content)

	return &dto.Message{
		ID:       id,
		SenderID: senderID,
		SendTime: time.Now(),
		Content:  content,
		Type:     1,
	}
}
