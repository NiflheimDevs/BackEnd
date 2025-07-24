package servicesimpl

import (
	"strconv"

	"github.com/niflheimdevs/backend/bootstrap"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/pkg"
)

type EmailService struct {
	Env      *bootstrap.Env
	UserRepo repositories.UserRepo
}

func NewEmailService(
	userRepo repositories.UserRepo,
	env *bootstrap.Env,
) *EmailService {
	return &EmailService{
		UserRepo: userRepo,
		Env:      env,
	}
}

func (es EmailService) SendNotification(userID int, subject string, body string) error {
	user, _ := es.UserRepo.FindUserByID(userID)

	if user.Is_verified {
		from := es.Env.Email.NotifAddr
		host := es.Env.Email.Host
		port, _ := strconv.Atoi(es.Env.Email.Port)
		username := es.Env.Email.Username
		password := es.Env.Email.Password
		message := pkg.NewMessage(from, user.Email, subject, body)
		err := pkg.SendEmail(host, port, username, password, message)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (es EmailService) SendEmail(to, subject string, body string) error {
	from := es.Env.Email.MainAddr
	host := es.Env.Email.Host
	port, _ := strconv.Atoi(es.Env.Email.Port)
	username := es.Env.Email.Username
	password := es.Env.Email.Password
	message := pkg.NewMessage(from, to, subject, body)
	err := pkg.SendEmail(host, port, username, password, message)
	if err != nil {
		return err
	}
	return nil
}
