package services

type NotifService interface {
	SendNotification(userID int, message string) error
}
