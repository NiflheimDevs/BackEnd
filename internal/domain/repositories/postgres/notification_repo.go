package repositories

import "github.com/niflheimdevs/backend/internal/domain/models"

type NotifRepo interface {
	SaveNotif(userID int, message string) int
	GetNotifs(userID int) []models.NotifModel
	GetNotif(id int) models.NotifModel
	ReadNotif(id int)
}
