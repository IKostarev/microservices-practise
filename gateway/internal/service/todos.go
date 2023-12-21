package service

import (
	"context"
	"fmt"
	"gateway/internal/models"
	"github.com/google/uuid"
)

func (s *GatewayService) CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	todo, err := s.todoServiceClient.CreateToDo(ctx, newTodo)
	if err != nil {
		return nil, fmt.Errorf("[CreateTodo service] create todo: %w", err)
	}

	return todo, nil
}

func (s *GatewayService) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	todo, err := s.todoServiceClient.UpdateToDo(ctx, newTodo)
	if err != nil {
		return nil, fmt.Errorf("[UpdateTodo service] upd todo: %w", err)
	}

	return todo, nil
}

func (s *GatewayService) GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error) {
	todoArr, err := s.todoServiceClient.GetToDos(ctx, todos)
	if err != nil {
		return nil, fmt.Errorf("[GetToDos service] get todos: %w", err)
	}

	return todoArr, nil
}

func (s *GatewayService) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error) {
	todo, err := s.todoServiceClient.GetToDo(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo service] get todo: %w", err)
	}

	return todo, nil
}

func (s *GatewayService) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	if err := s.todoServiceClient.DeleteToDo(ctx, todoID); err != nil {
		return fmt.Errorf("[DeleteToDo service] delete: %w", err)
	}

	return nil
}
