package grpc

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

func (g *GrpcApiServer) CreateToDo(ctx context.Context, newTodo *models.TodoDAO) (uuid.UUID, error) {
	todoID, err := g.todoService.CreateToDo(ctx, newTodo)
	if err != nil {
		g.logger.Error().Err(err).Msg("[CreateTodo] error")
		return uuid.Nil, err
	}

	return todoID, nil
}

func (g *GrpcApiServer) UpdateToDo(ctx context.Context, newTodo *models.TodoDAO) error {
	if err := g.todoService.UpdateToDo(ctx, newTodo); err != nil {
		g.logger.Error().Err(err).Msg("[UpdateToDo] error")
		return err
	}

	return nil
}

func (g *GrpcApiServer) GetToDos(ctx context.Context, todoID uuid.UUID) ([]models.TodoDAO, error) {
	todos, err := g.todoService.GetToDos(ctx, todoID)
	if err != nil {
		g.logger.Error().Err(err).Msg("[GetToDos] error")
		return nil, err
	}

	return todos, nil
}

func (g *GrpcApiServer) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error) {
	todo, err := g.todoService.GetToDo(ctx, todoID)
	if err != nil {
		g.logger.Error().Err(err).Msg("[GetToDo] error")
		return nil, err
	}

	return todo, nil
}

func (g *GrpcApiServer) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	if err := g.todoService.DeleteToDo(ctx, todoID); err != nil {
		g.logger.Error().Err(err).Msg("[DeleteToDo] error")
		return err
	}

	return nil
}
