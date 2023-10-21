package service

import (
	"context"
	"github.com/google/uuid"
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
	return nil, nil
}

func (s *TodoService) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	return nil, nil
}

func (s *TodoService) GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error) {
	return nil, nil
}

func (s *TodoService) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error) {
	return nil, nil
}

func (s *TodoService) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	return nil
}
