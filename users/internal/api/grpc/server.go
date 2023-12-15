package grpc

import (
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"users/config"
	"users/internal/api"
	"users/pkg/grpc_stubs/users"
)

type server struct {
	users.UnimplementedUserServiceServer
	userService api.UserService
}

func NewGrpcApi(
	cfg *config.Config,
	logger *zerolog.Logger,
	userService api.UserService,
) error {
	appAddr := fmt.Sprintf("%s:%s", cfg.Grpc.AppHost, cfg.Grpc.AppPort)
	lis, err := net.Listen("tcp", appAddr)
	if err != nil {
		return fmt.Errorf("[NewGrpcApi] listen: %w", err)
	}

	s := grpc.NewServer()

	logger.Info().Msgf("running GRPC server at '%s'", appAddr)
	users.RegisterUserServiceServer(s, &server{
		userService: userService,
	})
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("[NewGrpcApi] serve: %w", err)
	}

	return nil
}
