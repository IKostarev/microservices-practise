package config

import (
	"github.com/kelseyhightower/envconfig"
	"users/pkg/jaeger"
	"users/pkg/logging"
	"users/pkg/postgresql"
	"users/pkg/rabbitmq"
)

type Config struct {
	App           App                           `envconfig:"APP"`
	Grpc          Grpc                          `envconfig:"GRPC"`
	Password      PasswordConfig                `envconfig:"PASS"`
	Logging       logging.LoggerConfig          `envconfig:"LOG"`
	Postgres      postgresql.PostgreSQL         `envconfig:"POSTGRES"`
	Jaeger        jaeger.JaegerConfig           `envconfig:"JAEGER"`
	RabbitConfig  rabbitmq.RabbitProducerConfig `envconfig:"RABBITMQ"`
	UsersExchange string                        `envconfig:"RABBITMQ_USERS_EXCHANGE" default:"users.exchange"`
	UsersQueue    string                        `envconfig:"RABBITMQ_USERS_QUEUE" default:"users.queue"`
}

type MigrationsConfig struct {
	Postgres postgresql.PostgreSQL `envconfig:"POSTGRES"`
}

type App struct {
	AppHost string `envconfig:"APP_HOST" required:"true" default:"0.0.0.0"`
	AppPort string `envconfig:"APP_PORT" required:"true" default:"8000"`
}

type Grpc struct {
	AppHost string `envconfig:"GRPC_HOST" required:"true" default:"0.0.0.0"`
	AppPort string `envconfig:"GRPC_PORT" required:"true" default:"50000"`
}

type PasswordConfig struct {
	Time    uint32 `envconfig:"PASS_TIME" required:"true" default:"1"`
	Memory  uint32 `envconfig:"PASS_MEMORY" required:"true" default:"65536"`
	Threads uint8  `envconfig:"PASS_THREADS" required:"true" default:"4"`
	KeyLen  uint32 `envconfig:"PASS_KEY_LEN" required:"true" default:"32"`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess("", &c)
	return &c
}

func NewMigrationsFromEnv() *MigrationsConfig {
	c := MigrationsConfig{}
	envconfig.MustProcess("", &c)
	return &c
}
