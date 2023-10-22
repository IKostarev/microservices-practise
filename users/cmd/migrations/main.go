package main

import (
	"database/sql"
	"fmt"
	"users/pkg/postgresql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"users/config"
	"users/migrations"
)

func main() {
	cfg := config.NewMigrationsFromEnv()

	db, err := sql.Open("pgx", formDbURI(cfg.Postgres))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.MigrationsFS)
	err = goose.SetDialect("postgres")
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, ".", goose.WithAllowMissing())
	if err != nil {
		panic(err)
	}
}

func formDbURI(conf postgresql.PostgreSQL) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		conf.PostgresqlUser,
		conf.PostgresqlPassword,
		conf.PostgresqlHost,
		conf.PostgresqlPort,
		conf.PostgresqlDBName,
	)
}
