package app

import (
	"fmt"
	"gateway/config"
	api "gateway/internal/api"
	"gateway/internal/api/rest"
	"gateway/internal/clients/todos"
	"gateway/internal/clients/users"
	"gateway/internal/service"
	"golang.org/x/sync/errgroup"

	"gateway/pkg/logging"
	"github.com/rs/zerolog"
)

type App struct {
	cfg            *config.Config
	logger         *zerolog.Logger
	gatewayService api.GatewayService
}

func NewApp(cfg *config.Config) (*App, error) {
	logger := logging.NewLogger(cfg.Logging)

	todosClient, err := todos.NewTodosClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] grpc todo: %w", err)
	}

	usersClient, err := users.NewUsersClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] grpc users: %w", err)
	}

	gatewayService := service.NewGatewayService(&cfg.JWT, todosClient, usersClient)

	return &App{
		cfg:            cfg,
		logger:         logger,
		gatewayService: gatewayService,
	}, nil
}

func (a *App) RunAPI() error {
	group := new(errgroup.Group)

	group.Go(func() error {
		return rest.RunREST(a.cfg, a.logger, a.gatewayService)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunAPI] run: %w", err)
	}

	return nil
}
