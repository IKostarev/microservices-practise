package service

import (
	"context"
	"fmt"
	"gateway/internal/models"
	"github.com/opentracing/opentracing-go"
)

func (s *GatewayService) CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.CreateToDo")
	defer span.Finish()

	todoID, err := s.todoServiceClient.CreateToDo(ctx, newTodo)
	if err != nil {
		return 0, fmt.Errorf("[CreateToDo] store user:%w", err)
	}

	return todoID, nil
}

func (s *GatewayService) UpdateToDo(ctx context.Context, updateTodo *models.UpdateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdateTodo")
	defer span.Finish()

	todoID, err := s.todoServiceClient.UpdateToDo(ctx, updateTodo)
	if err != nil {
		return 0, fmt.Errorf("[UpdateToDo] update todo:%w", err)
	}

	return todoID, nil
}

func (s *GatewayService) GetToDos(ctx context.Context, todoID int) (*models.TodoDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetToDos")
	defer span.Finish()

	todo, err := s.todoServiceClient.GetToDos(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
	}

	return todo, nil
}

func (s *GatewayService) GetToDo(ctx context.Context, todoID int) (*models.TodoDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetTodo")
	defer span.Finish()

	todo, err := s.todoServiceClient.GetToDo(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
	}

	return todo, nil
}

func (s *GatewayService) DeleteToDo(ctx context.Context, todoID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.DeleteToDo")
	defer span.Finish()

	err := s.todoServiceClient.DeleteToDo(ctx, todoID)
	if err != nil {
		return fmt.Errorf("[DeleteToDo] delete todo:%w", err)
	}

	return nil
}
