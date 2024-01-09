package smtp_client

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
)

type SmtpClient struct {
	Config SmtpConfig
}

func NewSmtpClient(cfg SmtpConfig) *SmtpClient {
	return &SmtpClient{
		Config: cfg,
	}
}

func (s *SmtpClient) Send(to []string, subject, body string) error {
	// Sanitize input to prevent injection attacks
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	// Dial the SMTP server
	c, err := smtp.Dial(fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port))
	if err != nil {
		return err
	}
	defer c.Close()

	// Start TLS for secure connection
	if err = c.StartTLS(&tls.Config{ServerName: s.Config.Host}); err != nil {
		return err
	}

	// Authenticate with the server
	auth := smtp.PlainAuth("", s.Config.User, s.Config.Password, s.Config.Host)
	if err = c.Auth(auth); err != nil {
		return err
	}

	// Set the sender and recipient(s)
	if err = c.Mail(r.Replace(s.Config.User)); err != nil {
		return err
	}
	for _, recipient := range to {
		if err = c.Rcpt(r.Replace(recipient)); err != nil {
			return err
		}
	}

	// Send the email body
	w, err := c.Data()
	if err != nil {
		return err
	}

	// Compose the message
	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + s.Config.User + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	// Quit the session
	return c.Quit()
}
