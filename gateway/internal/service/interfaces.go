package service

import (
	"context"
	"gateway/internal/models"
)

type TodoServiceClient interface {
	CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error)
	UpdateToDo(ctx context.Context, updateTodo *models.UpdateTodoDTO) (int, error)
	GetToDos(ctx context.Context, todoID int) (*models.TodoDTO, error)
	GetToDo(ctx context.Context, todoID int) (*models.TodoDTO, error)
	DeleteToDo(ctx context.Context, todoID int) error
}

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
