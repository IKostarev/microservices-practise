package app

import (
	"context"
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
		return rabbitmq.ConsumeRabbitMessages(context.Background(), a.cfg, a.logger)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("[RunApp] run: %w", err)
	}

	return nil
}
