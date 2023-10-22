package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	maxConn           = 50
	healthCheckPeriod = 3 * time.Minute
	maxConnIdleTime   = 1 * time.Minute
	maxConnLifetime   = 3 * time.Minute
	minConns          = 10
	lazyConnect       = false
)

// NewPgxConn pool
func NewPgxConn(cfg *PostgreSQL) (*pgxpool.Pool, error) {
	// создадим DSN для подключения к базе
	ctx := context.Background()
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.PostgresqlHost,
		cfg.PostgresqlPort,
		cfg.PostgresqlUser,
		cfg.PostgresqlDBName,
		cfg.PostgresqlPassword,
	)

	// создадим конфиг пула соединений
	poolCfg, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}

	// установим конфигурацию соединения
	poolCfg.MaxConns = maxConn
	poolCfg.HealthCheckPeriod = healthCheckPeriod
	poolCfg.MaxConnIdleTime = maxConnIdleTime
	poolCfg.MaxConnLifetime = maxConnLifetime
	poolCfg.MinConns = minConns
	poolCfg.LazyConnect = lazyConnect

	// подключимся к базе данных
	connPool, err := pgxpool.ConnectConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	// вернем пул соединений
	return connPool, nil
}
