package service

import (
	"context"
	"fmt"
	"gateway/internal/app_errors"
	"gateway/internal/models"
	"gateway/pkg/ctxutil"
	"github.com/opentracing/opentracing-go"
)

func (s *GatewayService) RegisterUser(ctx context.Context, newUser *models.CreateUserDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.RegisterUser")
	defer span.Finish()

	// Передаем данные в слой репозитория для сохранения пользователя.
	userID, err := s.usersServiceClient.CreateUser(ctx, newUser)
	if err != nil {
		return 0, fmt.Errorf("[RegisterUser] store user:%w", err)
	}

	// возвращаем данные в слой хэндлера
	return userID, nil
}

func (s *GatewayService) UpdateUser(ctx context.Context, updatedUser *models.UserDTO) (*models.UserDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdateUser")
	defer span.Finish()

	senderID, ok := ctxutil.GetUserIDFromContext(ctx)
	if !ok {
		return nil, app_errors.ErrNoUserInContext
	}

	if senderID != updatedUser.ID {
		return nil, app_errors.NewUserIDMismatchError("UpdateUser", senderID, updatedUser.ID)
	}

	// Передаем данные в слой репозитория
	err := s.usersServiceClient.UpdateUser(ctx, updatedUser)
	if err != nil {
		return nil, fmt.Errorf("[UpdateUser] update user:%w", err)
	}

	// возвращаем данные в слой хэндлера
	return updatedUser, nil
}

func (s *GatewayService) UpdatePassword(ctx context.Context, request *models.UpdateUserPasswordDTO) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdatePassword")
	defer span.Finish()

	senderID, ok := ctxutil.GetUserIDFromContext(ctx)
	if !ok {
		return app_errors.ErrNoUserInContext
	}

	if senderID != request.ID {
		return app_errors.NewUserIDMismatchError("UpdatePassword", senderID, request.ID)
	}

	// Обновление пароля в базе данных
	err := s.usersServiceClient.UpdatePassword(ctx, request)
	if err != nil {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", err)
	}

	// возвращение ответа
	return nil
}

func (s *GatewayService) DeleteUser(ctx context.Context, userID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.DeleteUser")
	defer span.Finish()

	senderID, ok := ctxutil.GetUserIDFromContext(ctx)
	if !ok {
		return app_errors.ErrNoUserInContext
	}

	if senderID != userID {
		return app_errors.NewUserIDMismatchError("DeleteUser", senderID, userID)
	}

	// Удаление пользователя.
	err := s.usersServiceClient.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("[DeleteUser] delete user:%w", err)
	}

	return nil
}

func (s *GatewayService) GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetUserByID")
	defer span.Finish()

	// Получение пользователя по его идентификатору.
	storedUser, err := s.usersServiceClient.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByID] get user:%w", err)
	}

	// возврат данных пользователю
	return storedUser, nil
}

func (s *GatewayService) Login(ctx context.Context, login *models.UserLoginDTO) (*models.UserTokens, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Login")
	defer span.Finish()

	// Проверка наличия пользователя.
	existingUser, err := s.usersServiceClient.UserLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("[Login] log in: %w", err)
	}

	// Генерируем токены и возвращаем
	accessToken, err := s.jwtUtil.GenerateAccessToken(existingUser.ID)
	if err != nil {
		return nil, fmt.Errorf("[Login] generate access token:%w", err)
	}

	refreshToken, err := s.jwtUtil.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		return nil, fmt.Errorf("[Login] generate refresh token:%w", err)
	}

	return &models.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *GatewayService) Refresh(ctx context.Context, refresh string) (*models.UserTokens, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Refresh")
	defer span.Finish()

	userID, err := s.jwtUtil.VerifyToken(refresh)
	if err != nil {
		return nil, fmt.Errorf("[Refresh] verify token:%w", err)
	}

	// Генерируем токены и возвращаем
	accessToken, err := s.jwtUtil.GenerateAccessToken(userID)
	if err != nil {
		return nil, fmt.Errorf("[Refresh] generate access token:%w", err)
	}

	refreshToken, err := s.jwtUtil.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("[Refresh] generate refresh token:%w", err)
	}

	return &models.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
