package config

import (
	"github.com/kelseyhightower/envconfig"
	"notifications/pkg/logging"
	rabbitConfig "notifications/pkg/rabbitmq"
	"notifications/pkg/smtp_client"
)

type Config struct {
	App           App                       `envconfig:"APP"`
	Logging       logging.LoggerConfig      `envconfig:"LOG"`
	SmtpConfig    smtp_client.SmtpConfig    `envconfig:"SMTP"`
	RabbitConfig  rabbitConfig.RabbitConfig `envconfig:"RABBITMQ"`
	UsersExchange string                    `envconfig:"RABBITMQ_USERS_EXCHANGE" default:"users.exchange"`
	UsersQueue    string                    `envconfig:"RABBITMQ_USERS_QUEUE" default:"users.queue"`
}

type App struct {
	AppHost string `envconfig:"APP_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" required:"true" default:"8000"`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess("", &c)
	return &c
}
