package config

import (
	"github.com/kelseyhightower/envconfig"
	"users/pkg/jwtutil"
	"users/pkg/logging"
)

type Config struct {
	App      App                  `envconfig:"APP"`
	JWT      jwtutil.JWTUtil      `envconfig:"JWT"`
	Password PasswordConfig       `envconfig:"PASS"`
	Logging  logging.LoggerConfig `envconfig:"LOG"`
}

type App struct {
	AppHost string `envconfig:"APP_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" required:"true" default:"8000"`
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