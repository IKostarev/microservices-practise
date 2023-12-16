package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"notifications/internal/api"
	"notifications/internal/models"
	"notifications/pkg/rabbitmq/consumer"
)

type UsersMessages struct {
	usersService api.UserService
}

func NewUsersHandler(usersService api.UserService) *UsersMessages {
	return &UsersMessages{
		usersService: usersService,
	}
}

func (m *UsersMessages) Handle(
	ctx context.Context,
	msgList []amqp.Delivery,
) (uint64, []rabbitmq.Error, error) {
	var (
		lastSuccessConsumedTag uint64
		errorList              []rabbitmq.Error
	)

	for _, msg := range msgList {
		userMailItem, errUnmarshal := unmarshal(msg)
		if errUnmarshal != nil {
			errorList = append(errorList, *errUnmarshal)
			continue
		}

		err := m.usersService.SendUserMessage([]*models.UserMailItem{userMailItem})
		if err != nil {

		}
	}

	return lastSuccessConsumedTag, errorList, nil
}

func unmarshal(msg amqp.Delivery) (*models.UserMailItem, *rabbitmq.Error) {
	var item models.UserMailItem

	err := json.Unmarshal(msg.Body, &item)
	if err != nil {
		return nil, &rabbitmq.Error{
			Err: err,
			Msg: msg,
		}
	}

	return &item, nil
}
