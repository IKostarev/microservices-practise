package service

import (
	"context"
	"gateway/internal/models"
	"github.com/google/uuid"
)

type TodoServiceClient interface {
	CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error)
	UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error)
	GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error)
	GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error)
	DeleteToDo(ctx context.Context, todoID uuid.UUID) error
}

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/users_client.go -package=mocks gateway/internal/service UsersServiceClient
type UsersServiceClient interface {
	CreateUser(ctx context.Context, user *models.CreateUserDTO) (int, error)
	UpdateUser(ctx context.Context, user *models.UserDTO) error
	UpdatePassword(ctx context.Context, data *models.UpdateUserPasswordDTO) error
	DeleteUser(ctx context.Context, userID int) error
	GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error)
	GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*models.UserDTO, error)
	GetUserByUsername(ctx context.Context, username string) (*models.UserDTO, error)
	UserLogin(ctx context.Context, user *models.UserLoginDTO) (*models.UserDTO, error)
}
