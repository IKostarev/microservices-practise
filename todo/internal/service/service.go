package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
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

func (s *TodoService) CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	context, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	m := &models.TodoDAO{
		ID:          uuid.New(),
		CreatedBy:   rand.Intn(10),
		Assignee:    rand.Intn(15),
		Description: newTodo.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createResult, err := s.userRepo.CreateToDo(context, m)
	if err != nil {
		return nil, fmt.Errorf("[CreateToDo] error create is - %w\n", err)
	}

	return (*models.TodoDTO)(createResult), nil
}

func (s *TodoService) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	context, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	m := &models.TodoDAO{
		ID:          newTodo.ID,
		CreatedBy:   newTodo.CreatedBy,
		Assignee:    newTodo.Assignee,
		Description: newTodo.Description,
		CreatedAt:   newTodo.CreatedAt,
		UpdatedAt:   newTodo.UpdatedAt,
	}

	updateResult, err := s.userRepo.UpdateToDo(context, m)
	if err != nil {
		return nil, fmt.Errorf("[UpdateToDo] error update is - %w\n", err)
	}

	return (*models.TodoDTO)(updateResult), nil
}

func (s *TodoService) GetToDos(ctx context.Context) ([]models.TodoDTO, error) {
	context, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	list, err := s.userRepo.GetToDos(context)
	if err != nil {
		return nil, fmt.Errorf("[GetToDos] error get todos is - %w\n", err)
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
	context, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	todo, err := s.userRepo.GetToDo(context, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] error get todo is - %w", err)
	}

	return (*models.TodoDTO)(todo), nil
}

func (s *TodoService) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	context, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	if err := s.userRepo.DeleteToDo(context, todoID); err != nil {
		return fmt.Errorf("[DeleteToDo] error delete is - %w\n", err)
	}

	return nil
}
