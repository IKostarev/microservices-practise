package service

import (
	"context"
	"users/internal/models"
)

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/user_repository.go -package=mocks users/internal/service UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.CreateUserDTO) (int, error)
	UpdateUser(ctx context.Context, user *models.UserDAO) error
	UpdatePassword(ctx context.Context, userID int, newPassword string) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserByID(ctx context.Context, userID int) (*models.UserDAO, error)
	GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*models.UserDAO, error)
	GetUserByUsername(ctx context.Context, username string) (*models.UserDAO, error)
}

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/rabbit_producer.go -package=mocks users/internal/service RabbitProducer
type RabbitProducer interface {
	Publish(data []byte, reqesID string) (err error)
}

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/password_utils.go -package=mocks users/internal/service PasswordUtils
type PasswordUtils interface {
	GeneratePassword(ctx context.Context, password string) (string, error)
	ComparePassword(ctx context.Context, password, hash string) (bool, error)
}
