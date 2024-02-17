package grpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"todo/internal/models"
	"todo/pkg/ctxutil"
	"todo/pkg/grpc_stubs/todo"
)

func (s *server) CreateToDo(ctx context.Context, newTodo *todo.CreateTodoDTO) (*todo.TodoID, error) {
	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	newT := models.NewEmptyCreateTodoDTO().FromGRPC(newTodo)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.CreateToDo")
	defer span.Finish()

	id, err := s.todoService.CreateToDo(ctx, newT)

	requestId, _ := ctxutil.GetUserIDFromContext(ctx)
	s.logger.Info().Msgf("request id %s", requestId)

	if err != nil {
		s.logger.Error().
			Str("requestId", strconv.Itoa(requestId)).
			Msgf("[CreateToDo]: %w", err)
		return nil, err
	}

	return &todo.TodoID{Id: int32(id)}, nil
}

func (s *server) UpdateToDo(ctx context.Context, updateTodo *todo.UpdateTodoDTO) (*todo.TodoID, error) {
	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	upd := models.NewEmptyUpdateTodoDTO().FromGRPC(updateTodo)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.UpdateToDo")
	defer span.Finish()

	requestId, _ := ctxutil.GetUserIDFromContext(ctx)
	s.logger.Info().Msgf("request id %s", requestId)

	result, err := s.todoService.UpdateToDo(ctx, upd)
	if err != nil {
		s.logger.Error().
			Str("requestId", strconv.Itoa(requestId)).
			Msgf("[UpdateToDo]: %w", err)
		return nil, err
	}

	return &todo.TodoID{Id: int32(result)}, nil
}

func (s *server) GetToDos(ctx context.Context, getTodos *todo.TodoID) (*todo.TodoDTO, error) {
	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.UpdateToDo")
	defer span.Finish()

	requestId, _ := ctxutil.GetUserIDFromContext(ctx)
	s.logger.Info().Msgf("request id %s", requestId)

	_, err := s.todoService.GetToDo(ctx, int(getTodos.Id))
	if err != nil {
		s.logger.Error().
			Str("requestId", strconv.Itoa(requestId)).
			Msgf("[GetToDos]: %w", err)
		return nil, err
	}

	return nil, nil //TODO fix
}

func (s *server) GetToDo(ctx context.Context, getTodo *todo.TodoID) (*todo.TodoDTO, error) {
	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.UpdateToDo")
	defer span.Finish()

	requestId, _ := ctxutil.GetUserIDFromContext(ctx)
	s.logger.Info().Msgf("request id %s", requestId)

	_, err := s.todoService.GetToDo(ctx, int(getTodo.Id))
	if err != nil {
		s.logger.Error().
			Str("requestId", strconv.Itoa(requestId)).
			Msgf("[GetToDo]: %w", err)
		return nil, err
	}

	return nil, nil //TODO fix
}

func (s *server) DeleteToDo(ctx context.Context, todoID *todo.TodoID) (*emptypb.Empty, error) {
	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.DeleteToDo")
	defer span.Finish()

	requestId, _ := ctxutil.GetUserIDFromContext(ctx)
	s.logger.Info().Msgf("request id %s", requestId)

	err := s.todoService.DeleteToDo(ctx, int(todoID.Id))
	if err != nil {
		s.logger.Error().
			Str("requestId", strconv.Itoa(requestId)).
			Msgf("[DeleteToDo]: %w", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
