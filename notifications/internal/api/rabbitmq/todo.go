package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"notifications/internal/api"
	"notifications/internal/models"
)

type TodoMessagesHandler struct {
	logger      *zerolog.Logger
	todoService api.TodoService
}

func NewTodoMessagesHandler(
	logger *zerolog.Logger,
	todoService api.TodoService,
) *TodoMessagesHandler {
	return &TodoMessagesHandler{
		logger:      logger,
		todoService: todoService,
	}
}

func (m *TodoMessagesHandler) Handle(d amqp.Delivery) {
	requestID := d.Headers["requestId"]

	m.logger.Info().Msgf("request id %s", requestID)

	var item models.TodoMailItem
	err := json.Unmarshal(d.Body, &item)
	if err != nil {
		m.logger.Error().Msgf("[TodoMessagesHandler] unmarshalling message: %s", err)
		return
	}

	err = m.todoService.GetUserMessage(&item)
	if err != nil {
		m.logger.Error().Msgf("[TodoMessagesHandler] sending user mail: %s", err)
		return
	}

	d.Ack(true)
}
