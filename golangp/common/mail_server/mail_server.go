package mail_server

import (
	"bytes"
	"crypto/tls"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type EmailData struct {
	URL              string
	VerificationCode string
	FirstName        string
	Subject          string
	Amount           int64
	Balance          int64
}

type EmailConfig struct {
	EmailFrom string
	SMTPPass  string
	SMTPUser  string
	SMTPHost  string
	SMTPPort  int
}

// ? Email template parser

func ParseTemplateDir(dir string, defaultTemplatePath string) (*template.Template, error) {
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果目录不存在，使用默认路径
		dir = defaultTemplatePath
	}

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

func SendEmail(config *EmailConfig, email string, data *EmailData, emailTemp string, templatesPath string, defaultTemplatePath string) {
	// Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	tmpl, err := ParseTemplateDir(templatesPath, defaultTemplatePath)
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	err = tmpl.ExecuteTemplate(&body, emailTemp, &data)
	if err != nil {
		log.Fatal("Could not execute template verificationCode", err)
	}

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
		log.Fatal("Could not send email: ", err)
	}

}
