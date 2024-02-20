package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
	"todo/internal/models"
	"todo/pkg/redis"
)

const (
	GetTodoByIdCacheKey = "get todo by id %s"
)

type TodoRedisManager struct {
	redisManager *redis.RedisManager
	cacheTtl     time.Duration
}

func NewTodoRedisManager(redisManager *redis.RedisManager, cacheTtl time.Duration) *TodoRedisManager {
	return &TodoRedisManager{redisManager, cacheTtl}
}

func (m *TodoRedisManager) GetCacheByTodoID(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error) {
	var todo *models.TodoDAO

	redisTodo, err := m.redisManager.Client.Get(ctx, fmt.Sprintf(GetTodoByIdCacheKey, todoID.String())).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(redisTodo, &todo)
	if err != nil {
		m.FlushCache(ctx, todoID)

		return nil, err
	}

	return todo, nil
}

func (m *TodoRedisManager) StoreCache(ctx context.Context, todo *models.TodoDAO) {
	redisTodo, err := json.Marshal(todo)
	if err == nil {
		m.redisManager.Client.Set(
			ctx,
			fmt.Sprintf(GetTodoByIdCacheKey, todo.ID.String()),
			redisTodo,
			m.cacheTtl,
		)
	}
}

func (m *TodoRedisManager) FlushCache(ctx context.Context, todoID uuid.UUID) {
	m.redisManager.Client.Del(
		ctx,
		fmt.Sprintf(GetTodoByIdCacheKey, todoID.String()),
	)
}
