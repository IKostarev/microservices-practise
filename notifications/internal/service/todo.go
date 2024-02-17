package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"notifications/internal/models"
)

type TodoService struct {
	logger     *zerolog.Logger
	smtpClient SmtpClient
}

func NewTodoService(
	logger *zerolog.Logger,
	smtpClient SmtpClient,
) *TodoService {
	return &TodoService{
		logger:     logger,
		smtpClient: smtpClient,
	}
}

func (s *TodoService) GetUserMessage(item *models.TodoMailItem) error {

	var messageBody, subject string

	//switch {  TODO
	//default:
	//	return app_errors.ErrIncorrectUserEventType
	//}

	s.logger.Info().Msgf("[GetUserMessage] getting %s message", subject)
	err := s.smtpClient.Send(item.Receivers, subject, messageBody)
	if err != nil {
		return fmt.Errorf("[SendUserMessage]: send mssg: %w", err)
	}

	return nil
}
