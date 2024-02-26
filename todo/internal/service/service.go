package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"strconv"
	"todo/internal/models"
)

type TodoService struct {
	todoRepo           TodoRepository
	todoRabbitProducer RabbitProducer
	todoRedis          TodoRedisManager
}

func NewTodoService(
	todoRepo TodoRepository,
	todoRabbitProducer RabbitProducer,
	todoRedis TodoRedisManager,
) *TodoService {
	return &TodoService{
		todoRepo:           todoRepo,
		todoRabbitProducer: todoRabbitProducer,
		todoRedis:          todoRedis,
	}
}

func (t *TodoService) CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.CreateToDo")
	defer span.Finish()

	todoID, err := t.todoRepo.CreateToDo(ctx, newTodo)
	if err != nil {
		return 0, fmt.Errorf("[CreateToDo] store user:%w", err)
	}

	t.todoRedis.StoreCache(ctx, (*models.TodoDAO)(newTodo))

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

	extTodo, err := t.todoRedis.GetCacheByTodoID(ctx, updateTodo.ID)
	if err != nil {
		return 0, err
	}

	todoID, err := t.todoRepo.UpdateToDo(ctx, upd)
	if err != nil {
		return 0, fmt.Errorf("[UpdateToDo] update todo:%w", err)
	}

	t.todoRedis.StoreCache(ctx, extTodo)

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

	getTodo, err := t.todoRedis.GetCacheByTodoID(ctx, uuid.MustParse(strconv.Itoa(todoID)))
	if err != nil {
		getTodo, err = t.todoRepo.GetToDo(ctx, uuid.New())
		if err != nil {
			return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
		}
	}

	return (*models.TodoDTO)(getTodo), nil
}

func (t *TodoService) DeleteToDo(ctx context.Context, todoID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.DeleteToDo")
	defer span.Finish()

	_, err := t.todoRedis.GetCacheByTodoID(ctx, uuid.MustParse(strconv.Itoa(todoID)))
	if err != nil {
		err = t.todoRepo.DeleteToDo(ctx, uuid.New()) // TODO fix
		if err != nil {
			return fmt.Errorf("[DeleteToDo] delete todo:%w", err)
		}
	}

	t.todoRedis.FlushCache(ctx, uuid.MustParse(strconv.Itoa(todoID)))

	return nil
}
