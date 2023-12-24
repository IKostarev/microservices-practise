package grpc

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

func (g *GrpcApiServer) CreateToDo(ctx context.Context, newTodo *models.TodoDAO) (uuid.UUID, error) {
	dto := &models.TodoDTO{
		ID:          newTodo.ID,
		CreatedBy:   newTodo.CreatedBy,
		Assignee:    newTodo.Assignee,
		Description: newTodo.Description,
		CreatedAt:   newTodo.CreatedAt,
		UpdatedAt:   newTodo.UpdatedAt,
	}

	todoID, err := g.todoService.CreateToDo(ctx, dto)
	if err != nil {
		g.logger.Error().Err(err).Msg("[CreateTodo] error")
		return uuid.Nil, err
	}

	return todoID, nil
}

func (g *GrpcApiServer) UpdateToDo(ctx context.Context, newTodo *models.TodoDAO) error {
	dto := &models.TodoDTO{
		ID:          newTodo.ID,
		CreatedBy:   newTodo.CreatedBy,
		Assignee:    newTodo.Assignee,
		Description: newTodo.Description,
		CreatedAt:   newTodo.CreatedAt,
		UpdatedAt:   newTodo.UpdatedAt,
	}

	if err := g.todoService.UpdateToDo(ctx, dto); err != nil {
		g.logger.Error().Err(err).Msg("[UpdateToDo] error")
		return err
	}

	return nil
}

func (g *GrpcApiServer) GetToDos(ctx context.Context, todoID uuid.UUID) ([]models.TodoDAO, error) {
	list, err := g.todoService.GetToDos(ctx, todoID)
	if err != nil {
		g.logger.Error().Err(err).Msg("[GetToDos] error")
		return nil, err
	}

	result := make([]models.TodoDAO, 0, len(list))

	for _, item := range list {
		result = append(result, models.TodoDAO{
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

func (g *GrpcApiServer) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error) {
	todo, err := g.todoService.GetToDo(ctx, todoID)
	if err != nil {
		g.logger.Error().Err(err).Msg("[GetToDo] error")
		return nil, err
	}

	result := &models.TodoDAO{
		ID:          todo.ID,
		CreatedBy:   todo.CreatedBy,
		Assignee:    todo.Assignee,
		Description: todo.Description,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	return result, nil
}

func (g *GrpcApiServer) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	if err := g.todoService.DeleteToDo(ctx, todoID); err != nil {
		g.logger.Error().Err(err).Msg("[DeleteToDo] error")
		return err
	}

	return nil
}
