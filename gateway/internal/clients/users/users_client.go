package users

import (
	"context"
	"gateway/internal/models"
)

type UsersClient struct {
}

func NewUsersClient() *UsersClient {
	return &UsersClient{}
}

func (c *UsersClient) CreateUser(ctx context.Context, user *models.CreateUserDTO) (int, error) {
	return 0, nil
}
func (c *UsersClient) UpdateUser(ctx context.Context, user *models.UserDTO) error { return nil }
func (c *UsersClient) UpdatePassword(ctx context.Context, data *models.UpdateUserPasswordDTO) error {
	return nil
}
func (c *UsersClient) DeleteUser(ctx context.Context, userID int) error {
	return nil
}
func (c *UsersClient) GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	return nil, nil
}
func (c *UsersClient) GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*models.UserDTO, error) {
	return nil, nil
}
func (c *UsersClient) GetUserByUsername(ctx context.Context, username string) (*models.UserDTO, error) {
	return nil, nil
}
func (c *UsersClient) UserLogin(ctx context.Context, user *models.UserLoginDTO) (*models.UserDTO, error) {
	return nil, nil
}
