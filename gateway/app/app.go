package app

import (
	"fmt"
	"gateway/config"
	api "gateway/internal/api"
	"gateway/internal/api/rest"
	"gateway/internal/clients/todos"
	"gateway/internal/clients/users"
	"gateway/internal/service"
	"gateway/pkg/jaeger"
	"github.com/opentracing/opentracing-go"
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

	todosClient := todos.NewTodosClient()
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
	tracer, closer, err := jaeger.InitJaeger(&a.cfg.Jaeger, a.cfg.Logging.LogIndex)
	if err != nil {
		return fmt.Errorf("[NewApp] init jaeger %w", err)
	}

	a.logger.Info().Msgf("connected to jaeger at '%s'", a.cfg.Jaeger.Host)

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	group := new(errgroup.Group)

	group.Go(func() error {
		return rest.RunREST(a.cfg, a.logger, a.gatewayService)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunAPI] run: %w", err)
	}

	return nil
}
