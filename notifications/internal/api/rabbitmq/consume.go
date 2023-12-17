package rabbitmq

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"notifications/config"
	"notifications/internal/service"
	rabbitConsumer "notifications/pkg/rabbitmq/consumer"
	"notifications/pkg/smtp_client"
)

func ConsumeRabbitMessages(
	ctx context.Context,
	cfg *config.Config,
	logger *zerolog.Logger,
) error {
	logger.Info().Msg("starting notifications service")
	rabbitMQ, err := rabbitConsumer.New(&cfg.RabbitConfig, logger)
	if err != nil {
		return fmt.Errorf("[ConsumeRabbitMessages] connection to rabbitmq %w", err)
	}

	smtpClient := smtp_client.NewSmtpClient(cfg.SmtpConfig)
	usersService := service.NewUsersService(logger, smtpClient)
	usersMessagesHandler := NewUsersMessagesHandler(logger, usersService)

	rabbitMQ.SetHandler(
		cfg.UsersQueue,
		cfg.UsersExchange,
		usersMessagesHandler,
	)
	if err != nil {
		return fmt.Errorf("[ConsumeRabbitMessages] create users rabbitmq consumer: %w", err)
	}

	rabbitMQ.Run()

	<-ctx.Done()
	logger.Info().Msg("[ConsumeRabbitMessages] shutting down message consumption")

	return ctx.Err()
}
