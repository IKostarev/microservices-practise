package grpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"users/internal/models"
	"users/pkg/ctxutil"
	"users/pkg/grpc_stubs/users"
)

// RegisterUser - Handles user registration
func (s *server) RegisterUser(ctx context.Context, req *users.CreateUserDTO) (*users.UserID, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)
	newUser := models.NewEmptyCreateUserDTO().FromGRPC(req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.RegisterUser")
	defer span.Finish()

	id, err := s.userService.RegisterUser(ctx, newUser)

	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	s.logger.Info().Msgf("request id %s", requestId)

	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[RegisterUser]: %w", err)
		return nil, err
	}

	return &users.UserID{Id: int32(id)}, nil
}

// UpdateUser - Handles user updates
func (s *server) UpdateUser(ctx context.Context, req *users.UserDTO) (*users.UserDTO, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)
	updatedUser := models.NewEmptyUserDTO().FromGRPC(req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.UpdateUser")
	defer span.Finish()

	result, err := s.userService.UpdateUser(ctx, updatedUser)
	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdateUser]: %w", err)

		return nil, err
	}

	return result.ToGRPC(), nil
}

// UpdatePassword - Handles password updates
func (s *server) UpdatePassword(ctx context.Context, req *users.UpdateUserPasswordDTO) (*emptypb.Empty, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)
	updatePasswordDTO := models.NewEmptyUpdateUserPasswordDTO().FromGRPC(req)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.UpdatePassword")
	defer span.Finish()

	err := s.userService.UpdatePassword(ctx, updatePasswordDTO)
	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdatePassword]: %w", err)

		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteUser - Handles user deletion
func (s *server) DeleteUser(ctx context.Context, req *users.UserID) (*emptypb.Empty, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.DeleteUser")
	defer span.Finish()

	err := s.userService.DeleteUser(ctx, int(req.Id))
	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[DeleteUser]: %w", err)

		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetUserByID - Retrieves a user by ID
func (s *server) GetUserByID(ctx context.Context, req *users.UserID) (*users.UserDTO, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.GetUserByID")
	defer span.Finish()

	user, err := s.userService.GetUserByID(ctx, int(req.Id))
	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetUserByID]: %w", err)

		return nil, err
	}

	return user.ToGRPC(), nil
}

// GetUserByUsernameOrEmail - Retrieves a user by username or email
func (s *server) GetUserByUsernameOrEmail(ctx context.Context, req *users.UserDTO) (*users.UserDTO, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.GetUserByUsernameOrEmail")
	defer span.Finish()

	user, err := s.userService.GetUserByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetUserByUsernameOrEmail]: %w", err)

		return nil, err
	}

	return user.ToGRPC(), nil
}

// Login - Handles user login
func (s *server) Login(ctx context.Context, req *users.UserLoginDTO) (*users.UserDTO, error) {
	ctx = ctxutil.SetRequestIdFromGrpcToContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, "grpc_handler.Login")
	defer span.Finish()

	loginDTO := models.NewEmptyUserLoginDTO().FromGRPC(req)
	user, err := s.userService.Login(ctx, loginDTO)
	if err != nil {
		requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
		s.logger.Error().
			Str("requestId", requestId).
			Msgf("[Login]: %w", err)

		return nil, err
	}

	return user.ToGRPC(), nil
}
