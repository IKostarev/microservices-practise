package smtp

type SmtpConfig struct {
	Host     string `envconfig:"SMTP_HOST" default:"smtp.mail.ru"`
	Port     string `envconfig:"SMTP_PORT" default:"465"`
	User     string `envconfig:"SMTP_USER" default:"user"`
	Password string `envconfig:"SMTP_PASSWORD" default:"user"`
}
