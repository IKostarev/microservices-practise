package config

import (
	"github.com/kelseyhightower/envconfig"
	"todo/pkg/db/postgres"
	"todo/pkg/logger"
)

type App struct {
	AppHost string `envconfig:"APP_HOST" required:"true" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" required:"true" default:":8080"`
}

type Config struct {
	App        App                 `envconfig:"APP"`
	Logger     logger.LoggerConfig `envconfig:"LOGGER"`
	Database   postgres.PostgreSQL `envconfig:"POSTGRES"`
	Grpc       GrpcServerConfig    `envconfig:"GRPC"`
	Migrations MigrationsConfig
}

type MigrationsConfig struct {
	Postgres postgres.PostgreSQL `envconfig:"POSTGRES"`
}

type GrpcServerConfig struct {
	Host string `envconfig:"HOST" required:"true" required:"true" default:"0.0.0.0"`
	Port string `envconfig:"PORT" required:"true" required:"true" default:"50001"`
}

func LoadConfig() *Config {
	c := Config{}
	envconfig.MustProcess("", &c)
	c.Migrations = *LoadMigrationsConfig()
	return &c
}

func LoadMigrationsConfig() *MigrationsConfig {
	c := MigrationsConfig{}
	envconfig.MustProcess("", &c)
	return &c
}
