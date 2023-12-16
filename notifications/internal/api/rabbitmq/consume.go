package rabbitmq

import (
	"fmt"
	"github.com/rs/zerolog"
	"notifications/config"
	"notifications/internal/service"
	rabbitConsumer "notifications/pkg/rabbitmq/consumer"
	"notifications/pkg/smtp"
)

func ConsumeRabbitMessages(
	cfg *config.Config,
	logger *zerolog.Logger,
) error {
	rabbitMQ, err := rabbitConsumer.New(&cfg.RabbitConfig, logger)
	if err != nil {
		return fmt.Errorf("[ConsumeRabbitMessages] connection to rabbitmq %w", err)
	}

	smtpClient := smtp.NewSmtpClient(cfg.SmtpConfig)
	usersService := service.NewUsersService(smtpClient)
	usersMessagesHandler := NewUsersHandler(usersService)

	err = rabbitMQ.AddBatchConsumer(&cfg.RabbitConfig,
		cfg.UsersQueue,
		cfg.UsersExchange,
		usersMessagesHandler)
	if err != nil {
		return fmt.Errorf("[ConsumeRabbitMessages] create users rabbitmq consumer: %w", err)
	}

	return nil
}
