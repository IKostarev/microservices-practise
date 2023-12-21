package todos

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/models"
	"gateway/pkg/grpc_stubs/todo"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type TodosClient struct {
	client todo.TodoServiceClient
}

func NewTodosClient(cfg *config.Config, logger *zerolog.Logger) (*TodosClient, error) {
	appAdder := fmt.Sprintf("%s:%s", cfg.TodosClient.AppHost, cfg.TodosClient.AppGrpcPort)

	logger.Info().Msgf("[NewTodosClient] connecting via GRPC to todos at %s", appAdder)

	con, err := grpc.Dial(
		appAdder,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("[NewTodosClient] connect to todos GRPC service: %w", err)
	}

	c := todo.NewTodoServiceClient(con)
	return &TodosClient{client: c}, nil
}

func (c *TodosClient) CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	_, err := c.client.CreateTodo(ctx, newTodo.ToRPC())
	if err != nil {
		return nil, fmt.Errorf("[CreateToDo] create todo: %w", err)
	}

	return newTodo, nil

}
func (c *TodosClient) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	upd := &todo.TodoUpdate{
		Id:          newTodo.ID.String(),
		Description: newTodo.Description,
		Assignee:    int32(newTodo.Assignee),
		CreatedBy:   int32(newTodo.CreatedBy),
	}
	if _, err := c.client.UpdateTodo(ctx, upd); err != nil {
		return nil, fmt.Errorf("[UpdateToDo] update todo: %w", err)
	}

	return nil, nil // TODO fix
}
func (c *TodosClient) GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error) {
	res, err := c.client.GetTodos(ctx, &todo.TodoFilter{
		Assignee:  int32(todos.Assignee),
		CreatedBy: int32(todos.CreatedBy),
		DateFrom:  todos.DateFrom.String(),
		DateTo:    todos.DateTo.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("[GetToDos] get todos: %w", err)
	}

	todo := make([]models.TodoDTO, 0, len(res.Todos))
	for _, t := range res.Todos {
		todo = append(todo, models.TodoDTO{
			ID:          uuid.MustParse(t.Id),
			CreatedBy:   int(t.CreatedBy),
			Assignee:    int(t.Assignee),
			Description: t.Description,
			CreatedAt:   time.Unix(0, t.CreatedAt),
			UpdatedAt:   time.Unix(0, t.UpdatedAt),
		})
	}

	return todo, nil
}
func (c *TodosClient) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error) {
	res, err := c.client.GetTodo(ctx, &todo.TodoID{Id: todoID.String()})
	if err != nil {
		return nil, fmt.Errorf("[GetToDo] get todo: %w", err)
	}

	return &models.TodoDTO{
		ID:          todoID,
		CreatedBy:   int(res.CreatedBy),
		Assignee:    int(res.Assignee),
		Description: res.Description,
		CreatedAt:   time.Unix(0, res.CreatedAt),
		UpdatedAt:   time.Unix(0, res.UpdatedAt),
	}, nil
}
func (c *TodosClient) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	if _, err := c.client.DeleteTodo(ctx, &todo.TodoID{Id: todoID.String()}); err != nil {
		return fmt.Errorf("[DeleteToDo] delete todo: %w", err)
	}

	return nil
}
