package config

import (
	"github.com/kelseyhightower/envconfig"
	"todo/pkg/logger"
)

type App struct {
	AppHost string `envconfig:"APP_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" required:"true" default:":8080"`
}

type Config struct {
	App    App                 `envconfig:"APP"`
	Logger logger.LoggerConfig `envconfig:"LOGGER"`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess("", &c)
	return &c
}
