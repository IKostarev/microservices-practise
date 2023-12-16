package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"users/config"
	"users/internal/api"
	"users/internal/api/grpc"
	"users/internal/api/rest"
	"users/internal/repository"
	"users/internal/service"
	"users/pkg/logging"
	"users/pkg/postgresql"
	"users/pkg/rabbitmq/producer"
)

type App struct {
	cfg              *config.Config
	logger           *zerolog.Logger
	router           *mux.Router
	userService      api.UserService
	rabbitmqProducer *producer.Producer
}

func NewApp(
	cfg *config.Config,
) (*App, error) {
	logger := logging.NewLogger(cfg.Logging)

	// подключимся к базе данных
	databaseConn, err := postgresql.NewPgxConn(&cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// передадим подключение к базе данных констуктору репозитория
	userRepo := repository.NewUserRepository(databaseConn)

	// запустим rabbit mq продьюсер
	usersProducer, err := producer.New(
		&cfg.RabbitConfig,
		cfg.UsersExchange,
		cfg.UsersQueue,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("start rabbit producer: %w", err)
	}

	// передадим реализацию репозитория и продьюсера rabbit mq конструктору сервиса
	userService := service.NewUserService(&cfg.Password, userRepo, usersProducer)

	return &App{
		cfg:         cfg,
		logger:      logger,
		userService: userService,
	}, nil
}

func (a *App) RunApp() error {
	group := new(errgroup.Group)
	group.Go(func() error {
		err := rest.NewRestApi(a.cfg, a.logger, a.userService)
		return fmt.Errorf("[RunApp] run REST: %w", err)
	})

	group.Go(func() error {
		err := grpc.NewGrpcApi(a.cfg, a.logger, a.userService)
		return fmt.Errorf("[RunApp] run GRPC: %w", err)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunApp] run: %w", err)
	}

	return nil
}
