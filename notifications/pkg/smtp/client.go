package smtp

import (
	"fmt"
	smtpPkg "net/smtp"
)

type SmtpClient struct {
	Config SmtpConfig
}

func NewSmtpClient(cfg SmtpConfig) *SmtpClient {
	return &SmtpClient{
		Config: cfg,
	}
}

func (c *SmtpClient) Send(receivers []string, subject, body string) error {
	auth := smtpPkg.PlainAuth("", c.Config.User, c.Config.Password, c.Config.Host)
	err := smtpPkg.SendMail(
		fmt.Sprintf("%s:%s", c.Config.Host, c.Config.Port),
		auth,
		c.Config.User,
		receivers,
		[]byte(subject+body),
	)
	if err != nil {
		return err
	}

	return nil
}
