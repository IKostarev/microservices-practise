package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"users/internal/models"
	"users/pkg/grpc_stubs/users"
)

// RegisterUser - Handles user registration
func (s *server) RegisterUser(ctx context.Context, req *users.CreateUserDTO) (*users.UserID, error) {
	newUser := models.NewEmptyCreateUserDTO().FromGRPC(req)
	id, err := s.userService.RegisterUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &users.UserID{Id: int32(id)}, nil
}

// UpdateUser - Handles user updates
func (s *server) UpdateUser(ctx context.Context, req *users.UserDTO) (*users.UserDTO, error) {
	updatedUser := models.NewEmptyUserDTO().FromGRPC(req)
	result, err := s.userService.UpdateUser(ctx, updatedUser)
	if err != nil {
		return nil, err
	}
	return result.ToGRPC(), nil
}

// UpdatePassword - Handles password updates
func (s *server) UpdatePassword(ctx context.Context, req *users.UpdateUserPasswordDTO) (*emptypb.Empty, error) {
	updatePasswordDTO := models.NewEmptyUpdateUserPasswordDTO().FromGRPC(req)
	err := s.userService.UpdatePassword(ctx, updatePasswordDTO)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteUser - Handles user deletion
func (s *server) DeleteUser(ctx context.Context, req *users.UserID) (*emptypb.Empty, error) {
	err := s.userService.DeleteUser(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetUserByID - Retrieves a user by ID
func (s *server) GetUserByID(ctx context.Context, req *users.UserID) (*users.UserDTO, error) {
	user, err := s.userService.GetUserByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return user.ToGRPC(), nil
}

// GetUserByUsernameOrEmail - Retrieves a user by username or email
func (s *server) GetUserByUsernameOrEmail(ctx context.Context, req *users.UserDTO) (*users.UserDTO, error) {
	user, err := s.userService.GetUserByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	return user.ToGRPC(), nil
}

// Login - Handles user login
func (s *server) Login(ctx context.Context, req *users.UserLoginDTO) (*users.UserDTO, error) {
	loginDTO := models.NewEmptyUserLoginDTO().FromGRPC(req)
	user, err := s.userService.Login(ctx, loginDTO)
	if err != nil {
		return nil, err
	}
	return user.ToGRPC(), nil
}
