package servicesimpl

import (
	"encoding/json"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	websocket "github.com/niflheimdevs/backend/internal/infrastructure/websocket"
)

type NotifService struct {
	NotifRepo repositories.NotifRepo
	Hub       *websocket.Hub
}

func NewNotifService(
	notifRepo repositories.NotifRepo,
	hub *websocket.Hub,
) *NotifService {
	return &NotifService{
		NotifRepo: notifRepo,
		Hub:       hub,
	}
}

func (ns *NotifService) SendNotification(userID int, message string) error {
	id := ns.NotifRepo.SaveNotif(userID, message)

	notif := dto.NotifDTO{
		ID:       id,
		UserID:   userID,
		Content:  message,
		SendTime: time.Now(),
		IsRead:   false,
	}

	notifByte, err := json.Marshal(notif)

	if err != nil {
		return err
	}

	err = ns.Hub.SendToUser(userID, websocket.MessageTypeNotification, notifByte)

	if err != nil {
		return err
	}
	return nil
}
