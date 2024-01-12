package api

import (
	"context"
	"todo/internal/models"
)

type TodoService interface {
	CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error)
	UpdateToDo(ctx context.Context, updateTodo *models.UpdateTodoDTO) (int, error)
	GetToDos(ctx context.Context, todoID int) (*models.TodoDTO, error)
	GetToDo(ctx context.Context, todoID int) (*models.TodoDTO, error)
	DeleteToDo(ctx context.Context, todoID int) error
}
