package service

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

type TodoRepository interface {
	CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (uuid.UUID, error)
	UpdateToDo(ctx context.Context, updateTodo *models.TodoDAO) (uuid.UUID, error)
	GetToDos(ctx context.Context, todoID uuid.UUID) ([]models.TodoDAO, error)
	GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error)
	DeleteToDo(ctx context.Context, todoID uuid.UUID) error
}

type RabbitProducer interface {
	Publish(data []byte, reqesID string) (err error)
}
