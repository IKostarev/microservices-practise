package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
	"todo/pkg/jaeger"
	"todo/pkg/logging"
	"todo/pkg/postgresql"
	"todo/pkg/rabbitmq"
	"todo/pkg/redis"
)

type Config struct {
	App          App                           `envconfig:"APP"`
	Grpc         Grpc                          `envconfig:"GRPC"`
	Logging      logging.LoggerConfig          `envconfig:"LOG"`
	Postgres     postgresql.PostgreSQL         `envconfig:"POSTGRES"`
	Jaeger       jaeger.JaegerConfig           `envconfig:"JAEGER"`
	RabbitConfig rabbitmq.RabbitProducerConfig `envconfig:"RABBITMQ"`
	TodoExchange string                        `envconfig:"RABBITMQ_TODO_EXCHANGE" default:"todo.exchange"`
	TodoQueue    string                        `envconfig:"RABBITMQ_TODO_QUEUE" default:"todo.queue"`
	RedisConfig  redis.Config                  `envconfig:"REDIS"`
}

type MigrationsConfig struct {
	Postgres postgresql.PostgreSQL `envconfig:"POSTGRES"`
}

type App struct {
	AppHost     string        `envconfig:"APP_HOST" required:"true" default:"0.0.0.0"`
	AppPort     string        `envconfig:"APP_PORT" required:"true" default:"8000"`
	AppCacheTTL time.Duration `envconfig:"APP_CACHE_TTL" required:"true" default:"120m"`
}

type Grpc struct {
	AppHost string `envconfig:"GRPC_HOST" required:"true" default:"0.0.0.0"`
	AppPort string `envconfig:"GRPC_PORT" required:"true" default:"50000"`
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
