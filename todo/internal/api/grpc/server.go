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

type GrpcApiServer struct {
	todo.UnimplementedTodoServiceServer
	cfg         *config.Config
	todoService api.TodoService
	logger      *zerolog.Logger
}

func NewGrpcAPI(
	cfg *config.Config,
	logger *zerolog.Logger,
	todoService api.TodoService,
) *GrpcApiServer {
	return &GrpcApiServer{cfg: cfg, todoService: todoService, logger: logger}
}

func (s *GrpcApiServer) Start() error {
	srv := grpc.NewServer()
	todo.RegisterTodoServiceServer(srv, s)

	appAddr := fmt.Sprintf("%s:%s", s.cfg.Grpc.Host, s.cfg.Grpc.Port)
	lis, err := net.Listen("tcp", appAddr)
	if err != nil {
		return fmt.Errorf("[NewGrpcApi] listen: %w", err)
	}

	s.logger.Info().Msgf("running GRPC server at '%s'", appAddr)
	if err = srv.Serve(lis); err != nil {
		return fmt.Errorf("[NewGrpcApi] serve: %w", err)
	}

	return nil
}
