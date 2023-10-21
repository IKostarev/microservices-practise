package handlers

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

type TodoService interface {
	CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error)
	UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error)
	GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error)
	GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error)
	DeleteToDo(ctx context.Context, todoID uuid.UUID) error
}
