package pkg

import (
	"bytes"
	"html/template"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

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

func RenderTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func NewHTMLMessage(from, to, subject, htmlBody string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)
	msg.Embed(filepath.Join("internal", "application", "email_templates", "src", "BIDLANCERLOGO.svg"), gomail.SetHeader(map[string][]string{
		"Content-ID": {"<logo>"},
	}))
	return msg
}
