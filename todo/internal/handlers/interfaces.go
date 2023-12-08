package handlers

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

type TodoService interface {
	CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (uuid.UUID, error)
	UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) error
	GetToDos(ctx context.Context) ([]models.TodoDTO, error)
	GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error)
	DeleteToDo(ctx context.Context, todoID uuid.UUID) error
}
