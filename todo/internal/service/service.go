package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"todo/internal/models"
)

type TodoService struct {
	todoRepo           TodoRepository
	todoRabbitProducer RabbitProducer
}

func NewTodoService(
	todoRepo TodoRepository,
	todoRabbitProducer RabbitProducer,
) *TodoService {
	return &TodoService{
		todoRepo:           todoRepo,
		todoRabbitProducer: todoRabbitProducer,
	}
}

func (t *TodoService) CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.CreateToDo")
	defer span.Finish()

	todoID, err := t.todoRepo.CreateToDo(ctx, newTodo)
	if err != nil {
		return 0, fmt.Errorf("[CreateToDo] store user:%w", err)
	}

	return int(todoID.ID()), nil
}

func (t *TodoService) UpdateToDo(ctx context.Context, updateTodo *models.UpdateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdateTodo")
	defer span.Finish()

	upd := &models.TodoDAO{
		ID:          updateTodo.ID,
		Assignee:    updateTodo.Assignee,
		Description: updateTodo.Description,
		UpdatedAt:   updateTodo.UpdatedAt,
	}

	todoID, err := t.todoRepo.UpdateToDo(ctx, upd)
	if err != nil {
		return 0, fmt.Errorf("[UpdateToDo] update todo:%w", err)
	}

	return int(todoID.ID()), nil
}

func (t *TodoService) GetToDos(ctx context.Context, todoID int) (*models.TodoDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetToDos")
	defer span.Finish()

	_, err := t.todoRepo.GetToDos(ctx, uuid.New())
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
	}

	return nil, nil //TODO fix
}

func (t *TodoService) GetToDo(ctx context.Context, todoID int) (*models.TodoDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetTodo")
	defer span.Finish()

	todo, err := t.todoRepo.GetToDo(ctx, uuid.New())
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
	}

	return (*models.TodoDTO)(todo), nil
}

func (t *TodoService) DeleteToDo(ctx context.Context, todoID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.DeleteToDo")
	defer span.Finish()

	err := t.todoRepo.DeleteToDo(ctx, uuid.New()) // TODO fix
	if err != nil {
		return fmt.Errorf("[DeleteToDo] delete todo:%w", err)
	}

	return nil
}
