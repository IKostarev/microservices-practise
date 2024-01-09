package postgresql

import "time"

// PostgreSQL config.go
type PostgreSQL struct {
	PostgresqlHost     string        `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresqlPort     string        `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresqlUser     string        `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresqlPassword string        `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresqlDBName   string        `envconfig:"POSTGRES_DBNAME" default:"postgres"`
	MaxIdleConnTime    time.Duration `envconfig:"POSTGRES_MAX_IDLE_CONN_TIME" default:"5m"`
	MaxConns           int           `envconfig:"POSTGRES_MAX_CONNS" default:"20"`
	ConnMaxLifetime    time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"10m"`
}
