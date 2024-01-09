package service

type SmtpClient interface {
	Send(receivers []string, subject, body string) error
}
