package smtp

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/chheller/go-htmx-todo/modules/config"
)

// TODO: Refactor this to use a goroutine-safe singleton Client reference.
func SendEmail(to string, subject string, body string) error {
	env := config.GetEnvironment().SmtpConfig
	auth := smtp.PlainAuth("", env.Username, env.Password, env.Host)
	srvr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	from := fmt.Sprintf("%s < %s >", env.DisplayName, env.Username)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s", from, to, subject, mime, body))

	err := smtp.SendMail(srvr, auth, env.Username, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email %s", err)
		return err
	}
	return nil
}
