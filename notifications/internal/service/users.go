package service

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
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

func (s *UsersService) SendUserMessage(items []*models.UserMailItem) error {
	var errorsPool error

	for _, item := range items {
		var messageBody, subject string

		switch item.UserEventType {
		case models.UserEventTypeEmailVerification:
			messageBody = fmt.Sprintf(models.EmailBodyEmailVerification, item.Link)

		case models.UserEventTypeResetPassword:
		default:
			errorsPool = multierror.Append(errorsPool, app_errors.ErrIncorrectUserEventType)
			continue
		}

		err := s.smtpClient.Send(item.Receivers, subject, messageBody)
		if err != nil {
			errorsPool = multierror.Append(errorsPool, err)
		}
	}

	return errorsPool
}
