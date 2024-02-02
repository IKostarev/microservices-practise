package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	appErrors "users/internal/app_errors"
	"users/internal/models"
	"users/pkg/ctxutil"
)

type UserService struct {
	confirmationLink   string
	userRepo           UserRepository
	userRabbitProducer RabbitProducer
	passUtils          PasswordUtils
}

func NewUserService(
	userRepo UserRepository,
	userRabbitProducer RabbitProducer,
	utils PasswordUtils,
) *UserService {
	return &UserService{
		userRepo:           userRepo,
		userRabbitProducer: userRabbitProducer,
		passUtils:          utils,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, newUser *models.CreateUserDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.RegisterUser")
	defer span.Finish()

	// Проверка наличия пользователя с таким же именем или мэйлом.
	_, err := s.userRepo.GetUserByUsernameOrEmail(ctx, newUser.Username, newUser.Email)
	if err == nil {
		return 0, fmt.Errorf("[RegisterUser] get user: %w", appErrors.ErrUsernameOrEmailIsUsed)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	// проверяем, что пароль совпадает с подтверждением пароля
	if newUser.Password != newUser.PasswordConfirmation {
		return 0, fmt.Errorf("[RegisterUser] confirm pass: %w", appErrors.ErrPassAndConfirmationDoesNotMatch)
	}

	// Хеширование пароля - никогда не храните пароль в незашифрованном виде.
	hashedPassword, err := s.passUtils.GeneratePassword(ctx, newUser.Password)
	if err != nil {
		return 0, fmt.Errorf("[RegisterUser] generate pass: %w", err)
	}

	// укладываем хэш пароля вместо изначальноего представления
	newUser.Password = hashedPassword

	// Передаем данные в слой репозитория для сохранения пользователя.
	userID, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return 0, fmt.Errorf("[RegisterUser] store user:%w", err)
	}

	data, err := json.Marshal(models.UserMailItem{
		UserEventType: models.UserEventTypeEmailVerification,
		Receivers:     []string{newUser.Email},
		Link:          "example.com/verify",
	})
	if err != nil {
		return 0, fmt.Errorf("[RegisterUser] marshal email verification letter mssg:%w", err)
	}

	requestID, _ := ctxutil.GetRequestIDFromContext(ctx)
	err = s.userRabbitProducer.Publish(data, requestID)
	if err != nil {
		return 0, fmt.Errorf("[RegisterUser] publish email verification letter mssg:%w", err)
	}

	// возвращаем данные в слой хэндлера
	return userID, nil
}

func (s *UserService) UpdateUser(ctx context.Context, updatedUser *models.UserDTO) (*models.UserDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdateUser")
	defer span.Finish()

	// Проверка наличия пользователя.
	existingUser, err := s.userRepo.GetUserByID(ctx, updatedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("[UpdateUser] get user:%w", err)
	}

	// Обновление информации о пользователе.
	existingUser.Username = updatedUser.Username

	// Передаем данные в слой репозитория
	err = s.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		return nil, fmt.Errorf("[UpdateUser] update user:%w", err)
	}

	// возвращаем данные в слой хэндлера
	return updatedUser, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, request *models.UpdateUserPasswordDTO) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UpdatePassword")
	defer span.Finish()

	// Проверка наличия пользователя.
	existingUser, err := s.userRepo.GetUserByID(ctx, request.ID)
	if err != nil {
		return fmt.Errorf("[UpdatePassword] get user:%w", err)
	}

	// проверяем, совпадают ли новый пароль и подтверждение пароля
	if request.Password != request.PasswordConfirmation {
		return fmt.Errorf("[UpdatePassword] confirm pass:%w", appErrors.ErrPassAndConfirmationDoesNotMatch)
	}

	// Проверка старого пароля.
	passMatch, err := s.passUtils.ComparePassword(ctx, request.OldPassword, existingUser.Password)
	if err != nil {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", err)
	}
	if !passMatch {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", appErrors.ErrIncorrectOldPassword)
	}

	// Хеширование пароля.
	hashedPassword, err := s.passUtils.GeneratePassword(ctx, request.Password)
	if err != nil {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", err)
	}

	// Обновление пароля в базе данных
	err = s.userRepo.UpdatePassword(ctx, request.ID, hashedPassword)
	if err != nil {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", err)
	}

	// возвращение ответа
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.DeleteUser")
	defer span.Finish()

	// Удаление пользователя.
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("[DeleteUser] delete user:%w", err)
	}

	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetUserByID")
	defer span.Finish()

	// Получение пользователя по его идентификатору.
	var userResponse = new(models.UserDTO)
	storedUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return userResponse, fmt.Errorf("[GetUserByID] get user:%w", err)
	}

	// запись данных из DAO - data access object через которую мы работаем с базой данных
	// в DTO - data transfer object, который мы возвращаем пользователю
	userResponse.ID = storedUser.ID
	userResponse.Username = storedUser.Username
	userResponse.Email = storedUser.Email

	// возврат данных пользователю
	return userResponse, nil
}

func (s *UserService) GetUserByUsernameOrEmail(ctx context.Context, name, email string) (*models.UserDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetUserByUsernameOrEmail")
	defer span.Finish()

	// Получение пользователя по его идентификатору.
	var userResponse = new(models.UserDTO)
	storedUser, err := s.userRepo.GetUserByUsernameOrEmail(ctx, name, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return userResponse, fmt.Errorf("[GetUserByID] get user:%w", err)
	}

	// запись данных из DAO - data access object через которую мы работаем с базой данных
	// в DTO - data transfer object, который мы возвращаем пользователю
	userResponse.ID = storedUser.ID
	userResponse.Username = storedUser.Username
	userResponse.Email = storedUser.Email

	// возврат данных пользователю
	return userResponse, nil
}

func (s *UserService) Login(ctx context.Context, login *models.UserLoginDTO) (*models.UserDTO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Login")
	defer span.Finish()

	// Проверка наличия пользователя.
	existingUser, err := s.userRepo.GetUserByUsernameOrEmail(ctx, login.Username, login.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrWrongCredentials
		}
		return nil, fmt.Errorf("[Login] get user: %w", err)
	}

	// Проверка пароля
	passMatch, err := s.passUtils.ComparePassword(ctx, login.Password, existingUser.Password)
	if err != nil {
		return nil, fmt.Errorf("[Login] verify pass:%w", err)
	}
	if !passMatch {
		return nil, fmt.Errorf("[Login] verify pass:%w", appErrors.ErrWrongCredentials)
	}

	// возврат данных пользователю
	return &models.UserDTO{
		ID:       existingUser.ID,
		Username: existingUser.Username,
		Email:    existingUser.Email,
	}, nil
}
