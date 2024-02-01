package service

import (
	"gateway/pkg/jwtutil"
	"gateway/pkg/redis"
)

type GatewayService struct {
	jwtUtil            *jwtutil.JWTUtil
	todoServiceClient  TodoServiceClient
	usersServiceClient UsersServiceClient
	redisManager       *redis.RedisManager
}

func NewGatewayService(
	jwtUtil *jwtutil.JWTUtil,
	todoServiceClient TodoServiceClient,
	usersServiceClient UsersServiceClient,
	redisManager *redis.RedisManager,
) *GatewayService {
	return &GatewayService{
		jwtUtil:            jwtUtil,
		todoServiceClient:  todoServiceClient,
		usersServiceClient: usersServiceClient,
		redisManager:       redisManager,
	}
}
