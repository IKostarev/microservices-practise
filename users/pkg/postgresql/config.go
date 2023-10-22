package postgresql

import "time"

// PostgreSQL config
type PostgreSQL struct {
	PostgresqlHost     string        `envconfig:"HOST" default:"localhost"`
	PostgresqlPort     string        `envconfig:"PORT" yaml:"5432"`
	PostgresqlUser     string        `envconfig:"USER" yaml:"postgres"`
	PostgresqlPassword string        `envconfig:"PASSWORD" yaml:"postgres"`
	PostgresqlDBName   string        `envconfig:"DBNAME" yaml:"postgres"`
	MaxIdleConnTime    time.Duration `envconfig:"MAX_IDLE_CONN_TIME" default:"5m"`
	MaxConns           int           `envconfig:"MAX_CONNS" default:"20"`
	ConnMaxLifetime    time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"10m"`
}
