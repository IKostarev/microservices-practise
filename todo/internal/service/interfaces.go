package service

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

type TodoRepository interface {
	CreateToDo(ctx context.Context, newTodo *models.TodoDAO) (*models.TodoDAO, error)
	UpdateToDo(ctx context.Context, newTodo *models.TodoDAO) (*models.TodoDAO, error)
	GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDAO, error)
	GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error)
	DeleteToDo(ctx context.Context, todoID uuid.UUID) error
}
