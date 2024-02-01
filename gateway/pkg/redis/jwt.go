package redis

import (
	"context"
	"fmt"
	"time"
)

// IsTokenValid gets token from Redis by user ID and expiration time combined as a key.
// if not found in redis - token is invalid
func (rm *RedisManager) IsTokenValid(ctx context.Context, userID int, expiration int64) bool {
	result, err := rm.client.Exists(ctx, fmt.Sprintf("%d:%d", userID, expiration)).Result()
	if err != nil {
		return false
	}

	return result == 1
}

// StoreToken stores a token for a user in Redis with TTL.
func (rm *RedisManager) StoreToken(ctx context.Context, userID int, token string, expiration int64, ttlDuration time.Duration) error {
	err := rm.client.Set(ctx, fmt.Sprintf("%d:%d", userID, expiration), token, ttlDuration).Err()
	if err != nil {
		return fmt.Errorf("failed to store token in Redis: %v", err)
	}

	return nil
}

// InvalidateTokensForUser invalidates all tokens for a user.
func (rm *RedisManager) InvalidateTokensForUser(ctx context.Context, userID int) error {
	keys, err := rm.client.Keys(ctx, fmt.Sprintf("%d:*", userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys from Redis: %v", err)
	}

	for _, key := range keys {
		err := rm.client.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("failed to invalidate token in Redis: %v", err)
		}
	}

	return nil
}

// InvalidateToken invalidates a specific token for a user.
func (rm *RedisManager) InvalidateToken(ctx context.Context, userID int, expiration int64) error {
	err := rm.client.Del(ctx, fmt.Sprintf("%d:%d", userID, expiration)).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate token in Redis: %v", err)
	}

	return nil
}
