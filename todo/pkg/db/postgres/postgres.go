package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgreSQL struct {
	PostgresqlHost     string        `envconfig:"POSTGRES_HOST" default:"postgres"`
	PostgresqlPort     string        `envconfig:"POSTGRES_PORT" default:"5430"`
	PostgresqlUser     string        `envconfig:"POSTGRES_USER" default:"ibs_test"`
	PostgresqlPassword string        `envconfig:"POSTGRES_PASSWORD" default:"ibs_test"`
	PostgresqlDBName   string        `envconfig:"POSTGRES_DBNAME" default:"ibs_test"`
	MaxIdleConnTime    time.Duration `envconfig:"POSTGRES_MAX_IDLE_CONN_TIME" default:"5m"`
	MaxConns           int32         `envconfig:"POSTGRES_MAX_CONNS" default:"20"`
	ConnMaxLifetime    time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"10m"`
}

func NewPostgres(cfg *PostgreSQL) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresqlHost,
		cfg.PostgresqlPort,
		cfg.PostgresqlUser,
		cfg.PostgresqlPassword,
		cfg.PostgresqlDBName,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxIdleConnTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
