package handlers

import (
	"context"
	"users/internal/models"
)

type UserService interface {
	RegisterUser(ctx context.Context, newUser *models.CreateUserDTO) (int, error)
	UpdateUser(ctx context.Context, updatedUser *models.UserDTO) (*models.UserDTO, error)
	UpdatePassword(ctx context.Context, updatePassword *models.UpdateUserPasswordDTO) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error)

	Login(ctx context.Context, login *models.UserLoginDTO) (*models.UserTokens, error)
	Refresh(ctx context.Context, refresh string) (*models.UserTokens, error)
	VerifyToken(ctx context.Context, accessToken string) (int, error)
}
