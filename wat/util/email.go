package wat

import (
	"github.com/sirupsen/logrus"
	"net/smtp"
	"os"
	"regexp"
)

func IsValidEmail(email string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
func SendEmail(subject, content, address string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpHost := os.Getenv("SMTP_HOST")
	from := os.Getenv("SMTP_EMAIL")
	to := address
	password := os.Getenv("SMTP_PASSWORD")

	msg := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		content)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpServer, auth, from, []string{to}, msg)
	if err != nil {
		logrus.Errorf("Failed to send email: %v", err)
		return err
	}
	return nil
}
