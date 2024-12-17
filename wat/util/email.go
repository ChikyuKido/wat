package wat

import (
	"github.com/sirupsen/logrus"
	"net/smtp"
	"os"
)

func SendEmail(subject, content, address string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpHost := os.Getenv("SMTP_HOST")
	from := os.Getenv("SMTP_EMAIL")
	to := address
	password := os.Getenv("SMTP_PASSWORD")

	msg := []byte(subject + "\n" + content)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpServer, auth, from, []string{to}, msg)
	if err != nil {
		logrus.Errorf("Failed to send email: %v", err)
		return err
	}
	return nil
}
