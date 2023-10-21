package service

import (
	"context"
	"users/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.CreateUserDTO) (int, error)
	UpdateUser(ctx context.Context, user *models.UserDAO) error
	UpdatePassword(ctx context.Context, userID int, newPassword string) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserByID(ctx context.Context, userID int) (*models.UserDAO, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.UserDAO, error)
	GetUserByUsername(ctx context.Context, username string) (*models.UserDAO, error)
}
