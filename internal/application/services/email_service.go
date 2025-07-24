package services

type EmailService interface {
	SendNotification(userId int, subject string, body string) error
	SendEmail(to, subject string, body string) error
}
