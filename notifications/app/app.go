package app

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"notifications/config"
	"notifications/internal/api/rabbitmq"

	"github.com/rs/zerolog"
	"notifications/pkg/logging"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger
}

func NewApp(cfg *config.Config) (*App, error) {
	logger := logging.NewLogger(cfg.Logging)

	return &App{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (a *App) RunAPI() error {
	group := new(errgroup.Group)

	group.Go(func() error {
		err := rabbitmq.ConsumeRabbitMessages(a.cfg, a.logger)
		return fmt.Errorf("[RunApp] run rabbit consumer: %w", err)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunApp] run: %w", err)
	}

	return nil
}
