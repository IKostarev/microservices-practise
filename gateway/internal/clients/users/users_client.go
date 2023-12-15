package users

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/models"
	"gateway/pkg/grpc_stubs/users"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type UsersClient struct {
	client users.UserServiceClient
}

func NewUsersClient(cfg *config.Config, logger *zerolog.Logger) (*UsersClient, error) {
	appAddr := fmt.Sprintf("%s:%s", cfg.UsersClient.AppHost, cfg.UsersClient.AppGrpcPort)

	logger.Info().Msgf("[NewUsersClient] connecting via GRPC to usres at %s", appAddr)

	conn, err := grpc.Dial(
		appAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("[NewApp] connect to users GRPC service: %w", err)
	}

	c := users.NewUserServiceClient(conn)
	return &UsersClient{
		client: c,
	}, nil
}

func (c *UsersClient) CreateUser(ctx context.Context, user *models.CreateUserDTO) (int, error) {
	res, err := c.client.RegisterUser(ctx, user.ToGRPC())
	if err != nil {
		return 0, err
	}

	return int(res.Id), nil
}

func (c *UsersClient) UpdateUser(ctx context.Context, user *models.UserDTO) error {
	_, err := c.client.UpdateUser(ctx, &users.UserDTO{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *UsersClient) UpdatePassword(ctx context.Context, data *models.UpdateUserPasswordDTO) error {
	_, err := c.client.UpdatePassword(ctx, data.ToGRPC())
	if err != nil {
		return err
	}

	return nil
}

func (c *UsersClient) DeleteUser(ctx context.Context, userID int) error {
	_, err := c.client.DeleteUser(ctx, &users.UserID{
		Id: int32(userID),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *UsersClient) GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	user, err := c.client.GetUserByID(ctx, &users.UserID{
		Id: int32(userID),
	})
	if err != nil {
		return nil, err
	}

	return models.NewEmptyUserDTO().FromGRPC(user), nil
}

func (c *UsersClient) GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*models.UserDTO, error) {
	user, err := c.client.GetUserByUsernameOrEmail(ctx, &users.UserDTO{
		Username: username,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}

	return models.NewEmptyUserDTO().FromGRPC(user), nil
}

func (c *UsersClient) GetUserByUsername(ctx context.Context, username string) (*models.UserDTO, error) {
	user, err := c.client.GetUserByUsernameOrEmail(ctx, &users.UserDTO{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	return models.NewEmptyUserDTO().FromGRPC(user), nil
}

func (c *UsersClient) UserLogin(ctx context.Context, data *models.UserLoginDTO) (*models.UserDTO, error) {
	user, err := c.client.Login(ctx, data.ToGRPC())
	if err != nil {
		return nil, err
	}

	return models.NewEmptyUserDTO().FromGRPC(user), nil
}
