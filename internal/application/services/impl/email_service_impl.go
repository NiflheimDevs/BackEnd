package servicesimpl

import (
	"strconv"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/pkg"
)

type EmailService struct {
	Env *bootstrap.Env
}

func NewEmailService(env *bootstrap.Env) *EmailService {
	return &EmailService{
		Env: env,
	}
}

func (es EmailService) SendNotification(to, subject string, body string) error {
	from := es.Env.Email.NotifAddr
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
