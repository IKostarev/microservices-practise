package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"todo/config"
	"todo/internal/api"
	"todo/internal/api/grpc"
	"todo/internal/api/rest"
	"todo/internal/repository"
	"todo/internal/service"
	"todo/pkg/db/postgres"
	"todo/pkg/logger"
)

type App struct {
	cfg         *config.Config
	router      *mux.Router
	logger      *zerolog.Logger
	todoService api.TodoService
}

func NewApp(cfg *config.Config) (*App, error) {
	conn, err := postgres.NewPostgres(&cfg.Database)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	repo := repository.NewTodoRepository(conn)

	serv := service.NewTodoService(repo)

	log := logger.NewLogger(cfg.Logger)

	return &App{
		cfg:         cfg,
		logger:      log,
		todoService: serv,
	}, nil
}

func (a *App) RunApp() error {
	group := new(errgroup.Group)
	group.Go(func() error {
		err := rest.NewRestAPI(a.cfg, a.logger, a.todoService)
		return fmt.Errorf("[RunApp] run REST: %w", err)
	})

	group.Go(func() error {
		err := grpc.NewGrpcAPI(a.cfg, a.logger, a.todoService)
		return fmt.Errorf("[RunApp] run GRPC: %w", err)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunApp] run: %w", err)
	}

	return nil
}
