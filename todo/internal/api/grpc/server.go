package grpc

import (
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"todo/config"
	"todo/internal/api"
	"todo/pkg/grpc_stubs/todo"
)

type server struct {
	todo.UnimplementedTodoServiceServer
	todoService api.TodoService
	logger      *zerolog.Logger
}

func NewGrpcApi(
	cfg *config.Config,
	logger *zerolog.Logger,
	todoService api.TodoService,
) error {
	appAddr := fmt.Sprintf("%s:%s", cfg.Grpc.AppHost, cfg.Grpc.AppPort)
	lis, err := net.Listen("tcp", appAddr)
	if err != nil {
		return fmt.Errorf("[NewGrpcApi] listen: %w", err)
	}

	s := grpc.NewServer()

	logger.Info().Msgf("running GRPC server at '%s'", appAddr)
	todo.RegisterTodoServiceServer(s, &server{
		logger:      logger,
		todoService: todoService,
	})
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("[NewGrpcApi] serve: %w", err)
	}

	return nil
}
