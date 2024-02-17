package rest

import (
	"context"
	"gateway/internal/models"
)

type GatewayService interface {
	RegisterUser(ctx context.Context, newUser *models.CreateUserDTO) (int, error)
	UpdateUser(ctx context.Context, updatedUser *models.UserDTO) (*models.UserDTO, error)
	UpdatePassword(ctx context.Context, updatePassword *models.UpdateUserPasswordDTO) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error)
	Login(ctx context.Context, login *models.UserLoginDTO) (*models.UserTokens, error)
	Refresh(ctx context.Context, refresh string, access string) (*models.UserTokens, error)
	InvalidateTokensForUser(ctx context.Context, userID int) error
	InvalidateToken(ctx context.Context, userID int, access, refresh string) error

	CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error)
	UpdateToDo(ctx context.Context, updateTodo *models.UpdateTodoDTO) (int, error)
	GetToDos(ctx context.Context, todoID int) (*models.TodoDTO, error)
	GetToDo(ctx context.Context, todoID int) (*models.TodoDTO, error)
	DeleteToDo(ctx context.Context, todoID int) error
}
