package service

import (
	"fmt"
	"notifications/internal/app_errors"
	"notifications/internal/models"
)

type UsersService struct {
	smtpClient SmtpClient
}

func NewUsersService(smtpClient SmtpClient) *UsersService {
	return &UsersService{
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

	err := s.smtpClient.Send(item.Receivers, subject, messageBody)
	if err != nil {
		return err
	}

	return nil
}
