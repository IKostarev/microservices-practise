package service

import (
	"fmt"
	"github.com/rs/zerolog"
	"notifications/internal/app_errors"
	"notifications/internal/models"
)

type UsersService struct {
	logger     *zerolog.Logger
	smtpClient SmtpClient
}

func NewUsersService(
	logger *zerolog.Logger,
	smtpClient SmtpClient,
) *UsersService {
	return &UsersService{
		logger:     logger,
		smtpClient: smtpClient,
	}
}

func (s *UsersService) SendUserMessage(item *models.UserMailItem) error {

	var messageBody, subject string

	switch item.UserEventType {
	case models.UserEventTypeEmailVerification:
		messageBody = fmt.Sprintf(models.EmailBodyEmailVerification, item.Link)
		subject = models.EmailSubjectEmailVerification

	default:
		return app_errors.ErrIncorrectUserEventType
	}

	s.logger.Info().Msgf("[SendUserMessage] sending %s message", subject)
	err := s.smtpClient.Send(item.Receivers, subject, messageBody)
	if err != nil {
		return fmt.Errorf("[SendUserMessage]: send mssg: %w", err)
	}

	return nil
}
