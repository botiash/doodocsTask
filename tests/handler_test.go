package tests

import (
	"net/smtp"
	"testing"

	"doodocsProg/internal/config"
)

func TestGetArchiveInfoHandler(t *testing.T) {
}

func TestCreateArchiveHandler(t *testing.T) {
}
func sendTestEmail() error {
	cnf := config.Get()
	auth := smtp.PlainAuth("", cnf.Username, cnf.Password, cnf.Smtp_server)
	from := ""
	to := []string{"akylbek.ba@mail.ru"}
	message := []byte("To: akylbek.ba@mail.ru\r\n" +
		"Subject: Test Email\r\n" +
		"\r\n" +
		"This is a test email.")

	err := smtp.SendMail(cnf.Smtp_server, auth, from, to, message)
	return err
}

func TestSendEmail(t *testing.T) {
	err := sendTestEmail()
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
		return
	}
	t.Log("Email sent successfully")
}
