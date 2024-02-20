package service

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/models"
)

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/todo_repository.go -package=mocks todo/internal/service TodoRepository
type TodoRepository interface {
	CreateToDo(ctx context.Context, newTodo *models.CreateTodoDTO) (uuid.UUID, error)
	UpdateToDo(ctx context.Context, updateTodo *models.TodoDAO) (uuid.UUID, error)
	GetToDos(ctx context.Context, todoID uuid.UUID) ([]models.TodoDAO, error)
	GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error)
	DeleteToDo(ctx context.Context, todoID uuid.UUID) error
}

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/rabbit_producer.go -package=mocks todo/internal/service RabbitProducer
type RabbitProducer interface {
	Publish(data []byte, reqesID string) (err error)
}

//go:generate mockgen --build_flags=-mod=mod -destination=./mocks/todo_redis_manager.go -package=mocks todo/internal/service TodoRedisManager
type TodoRedisManager interface {
	GetCacheByTodoID(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error)
	StoreCache(ctx context.Context, todo *models.TodoDAO)
	FlushCache(ctx context.Context, todoID uuid.UUID)
}
