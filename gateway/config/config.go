package config

import (
	"gateway/pkg/jwtutil"
	"gateway/pkg/logging"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App         App                  `envconfig:"APP"`
	JWT         jwtutil.JWTUtil      `envconfig:"JWT"`
	Logging     logging.LoggerConfig `envconfig:"LOG"`
	UsersClient UsersClient          `envconfig:"USERS"`
	TodosClient TodosClient          `envconfig:"TODOS"`
}

type App struct {
	AppHost string `envconfig:"APP_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" required:"true" default:"8000"`
}

type UsersClient struct {
	AppHost string `envconfig:"USERS_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"USERS_PORT" required:"true" default:"8000"`
}

type TodosClient struct {
	AppHost string `envconfig:"TODOS_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"TODOS_PORT" required:"true" default:"8000"`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess("", &c)
	return &c
}
