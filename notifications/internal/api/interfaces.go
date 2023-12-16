package api

import "notifications/internal/models"

type UserService interface {
	SendUserMessage(items []*models.UserMailItem) error
}
