package smtp

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/chheller/go-htmx-todo/modules/config"
)

// TODO: Refactor this to use a goroutine-safe singleton Client reference.
func SendEmail(to string, subject string, body string) {
	env := config.GetEnvironment().SmtpConfig
	auth := smtp.PlainAuth("", env.Username, env.Password, env.Host)
	srvr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s", to, subject, body))

	err := smtp.SendMail(srvr, auth, env.Username, []string{to}, msg)
	if err != nil {
		log.Panicf("Error sending email %s", err)
	}
}
