package todos

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/models"
	"gateway/pkg/ctxutil"
	"gateway/pkg/grpc_stubs/todo"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type TodosClient struct {
	client todo.TodoServiceClient
}

func NewTodosClient(cfg *config.Config, logger *zerolog.Logger) (*TodosClient, error) {
	appAddr := fmt.Sprintf("%s:%s", cfg.UsersClient.AppHost, cfg.TodosClient.AppGrpcPort)

	logger.Info().Msgf("[NewTodosClient] connecting via GRPC to todo at %s", appAddr)

	conn, err := grpc.Dial(
		appAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("[NewTodosClient] connect to todo GRPC service: %w", err)
	}

	c := todo.NewTodoServiceClient(conn)
	return &TodosClient{
		client: c,
	}, nil
}

func (c *TodosClient) CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client.CreateToDo")
	defer span.Finish()

	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	todoID, err := c.CreateToDo(ctx, newTodo)
	if err != nil {
		return 0, fmt.Errorf("[CreateToDo] store user:%w", err)
	}

	return todoID, nil
}

func (c *TodosClient) UpdateToDo(ctx context.Context, updateTodo *models.UpdateTodoDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client.UpdateTodo")
	defer span.Finish()

	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	todoID, err := c.UpdateToDo(ctx, updateTodo)
	if err != nil {
		return 0, fmt.Errorf("[UpdateToDo] update todo:%w", err)
	}

	return todoID, nil
}

func (c *TodosClient) GetToDos(ctx context.Context, todoID int) (*models.TodoDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client.GetToDos")
	defer span.Finish()

	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	todo, err := c.GetToDos(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
	}

	return todo, nil
}

func (c *TodosClient) GetToDo(ctx context.Context, todoID int) (*models.TodoDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client.GetTodo")
	defer span.Finish()

	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	todo, err := c.GetToDo(ctx, todoID)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo:%w", err)
	}

	return todo, nil
}

func (c *TodosClient) DeleteToDo(ctx context.Context, todoID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client.DeleteToDo")
	defer span.Finish()

	ctx = ctxutil.SetRequestIdFromContextToGrpc(ctx)
	err := c.DeleteToDo(ctx, todoID)
	if err != nil {
		return fmt.Errorf("[DeleteToDo] delete todo:%w", err)
	}

	return nil
}
