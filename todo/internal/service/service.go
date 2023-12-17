package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"todo/internal/models"
)

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type TodoService struct {
	passConfig *PasswordConfig
	userRepo   TodoRepository
}

func NewTodoService(userRepo TodoRepository) *TodoService {
	return &TodoService{
		userRepo: userRepo,
	}
}

func (s *TodoService) CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (uuid.UUID, error) {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	m := &models.TodoDAO{
		ID:          uuid.New(),
		CreatedBy:   newTodo.CreatedBy,
		Assignee:    newTodo.Assignee,
		Description: newTodo.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createResult, err := s.userRepo.CreateToDo(context, m)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[CreateToDo] create - %w\n", err)
	}

	return createResult, nil
}

func (s *TodoService) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) error {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	m := &models.TodoDAO{
		ID:          newTodo.ID,
		CreatedBy:   newTodo.CreatedBy,
		Assignee:    newTodo.Assignee,
		Description: newTodo.Description,
		CreatedAt:   newTodo.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	err := s.userRepo.UpdateToDo(context, m)
	if err != nil {
		return fmt.Errorf("[UpdateToDo] update - %w\n", err)
	}

	return nil
}

func (s *TodoService) GetToDos(ctx context.Context, todoID uuid.UUID) ([]models.TodoDTO, error) {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	list, err := s.userRepo.GetToDos(context, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDos] get todos - %w\n", err)
	}

	result := make([]models.TodoDTO, 0, len(list))

	for _, item := range list {
		result = append(result, models.TodoDTO{
			ID:          item.ID,
			CreatedBy:   item.CreatedBy,
			Assignee:    item.Assignee,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return result, nil
}

func (s *TodoService) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error) {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	todo, err := s.userRepo.GetToDo(context, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo - %w", err)
	}

	return (*models.TodoDTO)(todo), nil
}

func (s *TodoService) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	if err := s.userRepo.DeleteToDo(context, todoID); err != nil {
		return fmt.Errorf("[DeleteToDo] delete - %w\n", err)
	}

	return nil
}
