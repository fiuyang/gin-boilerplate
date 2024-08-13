package utils

import (
	"bytes"
	"crypto/tls"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
	"scylla/entity"
	"scylla/pkg/config"
	"scylla/pkg/exception"
)

type EmailData struct {
	Otp     int
	Email   string
	Subject string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *entity.User, data *EmailData, emailTemp string) {
	conf := config.Get()
	// Sender data.
	from := conf.Email.From
	smtpPass := conf.Email.Pass
	smtpUser := conf.Email.User
	to := user.Email
	smtpHost := conf.Email.Host
	smtpPort := conf.Email.Port

	var body bytes.Buffer

	template, err := ParseTemplateDir("template")
	if err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

	template.ExecuteTemplate(&body, emailTemp, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		panic(exception.NewInternalServerErrorHandler(err.Error()))
	}

}
