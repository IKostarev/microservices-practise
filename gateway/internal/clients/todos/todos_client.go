package todos

import (
	"context"
	"gateway/internal/models"
	"github.com/google/uuid"
)

type TodosClient struct {
}

func NewTodosClient() *TodosClient {
	return &TodosClient{}
}

func (c *TodosClient) CreateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	return nil, nil
}
func (c *TodosClient) UpdateToDo(ctx context.Context, newTodo *models.TodoDTO) (*models.TodoDTO, error) {
	return nil, nil
}
func (c *TodosClient) GetToDos(ctx context.Context, todos *models.GetTodosDTO) ([]models.TodoDTO, error) {
	return nil, nil
}
func (c *TodosClient) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDTO, error) {
	return nil, nil
}
func (c *TodosClient) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	return nil
}
