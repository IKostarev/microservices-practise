package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"todo/config"
	"todo/internal/api"
	"todo/internal/api/grpc"
	"todo/internal/api/rest"
	"todo/internal/repository"
	"todo/internal/service"
	"todo/pkg/jaeger"
	"todo/pkg/logging"
	"todo/pkg/postgresql"
	"todo/pkg/rabbitmq/producer"
)

type App struct {
	cfg              *config.Config
	router           *mux.Router
	logger           *zerolog.Logger
	todoService      api.TodoService
	rabbitmqProducer *producer.Producer
}

func NewApp(
	cfg *config.Config,
) (*App, error) {
	logger := logging.NewLogger(cfg.Logging)

	databaseConn, err := postgresql.NewPgxConn(&cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	todoRepo := repository.NewTodoRepository(databaseConn)

	todosProducer, err := producer.New(
		&cfg.RabbitConfig,
		cfg.TodoExchange,
		cfg.TodoQueue,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("start rabbit producer: %w", err)
	}

	todoService := service.NewTodoService(todoRepo, todosProducer)

	return &App{
		cfg:         cfg,
		logger:      logger,
		todoService: todoService,
	}, nil
}

func (a *App) RunApp() error {
	tracer, closer, err := jaeger.InitJaeger(&a.cfg.Jaeger, a.cfg.Logging.LogIndex)
	if err != nil {
		return fmt.Errorf("[RunApp] init jaeger %w", err)
	}

	a.logger.Info().Msgf("connect to jaeger at '%s'", a.cfg.Jaeger.Host)

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	group := new(errgroup.Group)
	group.Go(func() error {
		err := rest.NewRestApi(a.cfg, a.logger, a.todoService)
		return fmt.Errorf("[RunApp] run REST: %w", err)
	})

	group.Go(func() error {
		err := grpc.NewGrpcApi(a.cfg, a.logger, a.todoService)
		return fmt.Errorf("[RunApp] run GRPC: %w", err)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunApp] run: %w", err)
	}

	return nil
}
