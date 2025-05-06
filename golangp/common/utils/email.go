package utils

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"github.com/toheart/functrace"
	"gopkg.in/gomail.v2"
)

type AccountEmailData struct {
	URL              string
	VerificationCode string
	FirstName        string
	Subject          string
	Amount           int64
	Balance          int64
}

type ShareProjectEmailData struct {
	AuthorizedUser string
	SharerUser     string
	SharerEmail    string
	ProjectName    string
	ProjectLink    string
	Subject        string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	defer functrace.Trace([]interface {
	}{dir})()
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

func SendEmail(emailTo string, subject string, data any, emailTemp string) {
	defer functrace.Trace([]interface {
	}{emailTo, subject, data, emailTemp})()
	configCopy := config.Env

	from := "Notification <" + configCopy.EmailFrom + ">"
	smtpPass := configCopy.SMTPPass
	smtpUser := configCopy.SMTPUser
	to := emailTo
	smtpHost := configCopy.SMTPHost
	smtpPort := configCopy.SMTPPort

	var body bytes.Buffer

	tmpl, err := ParseTemplateDir("common/templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
		return
	}

	err = tmpl.ExecuteTemplate(&body, emailTemp, &data)
	if err != nil {
		log.Fatalf("Could not execute template %s, %v", emailTemp, err)
		return
	}

	htmlContent := body.String()

	plainContent := html2text.HTML2Text(htmlContent)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", plainContent)
	m.AddAlternative("text/html", htmlContent)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}

func SendAccountEmail(user *models.User, data *AccountEmailData, emailTemp string) {
	defer functrace.Trace([]interface {
	}{user, data, emailTemp})()
	SendEmail(user.Email, data.Subject, data, emailTemp)
}

func SendShareProjectEmail(user *models.User, data *ShareProjectEmailData, emailTemp string) {
	defer functrace.Trace([]interface {
	}{user, data, emailTemp})()
	SendEmail(user.Email, data.Subject, data, emailTemp)
}
