package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"notifications/internal/api"
	"notifications/internal/models"
)

type UsersMessagesHandler struct {
	logger       *zerolog.Logger
	usersService api.UserService
}

func NewUsersMessagesHandler(
	logger *zerolog.Logger,
	usersService api.UserService,
) *UsersMessagesHandler {
	return &UsersMessagesHandler{
		logger:       logger,
		usersService: usersService,
	}
}

func (m *UsersMessagesHandler) Handle(d amqp.Delivery) {
	var item models.UserMailItem
	err := json.Unmarshal(d.Body, &item)
	if err != nil {
		m.logger.Error().Msgf("[UsersMessagesHandler] unmarshalling message: %s", err)
		return
	}

	err = m.usersService.SendUserMessage(&item)
	if err != nil {
		m.logger.Error().Msgf("[UsersMessagesHandler] sending user mail: %s", err)
		return
	}

	d.Ack(true)
}
