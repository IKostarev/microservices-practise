package service

import (
	"gateway/pkg/jwtutil"
)

type GatewayService struct {
	jwtUtil            *jwtutil.JWTUtil
	todoServiceClient  TodoServiceClient
	usersServiceClient UsersServiceClient
}

func NewGatewayService(
	jwtUtil *jwtutil.JWTUtil,
	todoServiceClient TodoServiceClient,
	usersServiceClient UsersServiceClient,
) *GatewayService {
	return &GatewayService{
		jwtUtil:            jwtUtil,
		todoServiceClient:  todoServiceClient,
		usersServiceClient: usersServiceClient,
	}
}
