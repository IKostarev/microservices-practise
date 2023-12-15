package service

import (
	"context"
	"gateway/internal/models"
	"github.com/google/uuid"
)

func (s *GatewayService) CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	return nil, nil
}

func (s *GatewayService) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	return nil, nil
}

func (s *GatewayService) GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error) {
	return nil, nil
}

func (s *GatewayService) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error) {
	return nil, nil
}

func (s *GatewayService) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	return nil
}
