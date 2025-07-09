package pkg

import "gopkg.in/gomail.v2"

func NewMessage(from, to, subject string, body string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)
	return msg
}

func SendEmail(host string, port int, username, password string, msg *gomail.Message) error {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.SSL = true

	err := dialer.DialAndSend(msg)
	return err
}
