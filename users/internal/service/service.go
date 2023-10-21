package service

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"users/config"
	"users/internal/models"
	"users/pkg/jwtutil"

	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

type UserService struct {
	passConfig *config.PasswordConfig
	userRepo   UserRepository
	jwtUtil    *jwtutil.JWTUtil
}

func NewUserService(
	passwordConfig *config.PasswordConfig,
	userRepo UserRepository,
	util *jwtutil.JWTUtil,
) *UserService {
	return &UserService{
		passConfig: passwordConfig,
		userRepo:   userRepo,
		jwtUtil:    util,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, newUser *models.CreateUserDTO) (int, error) {
	// Проверка наличия пользователя с таким же именем или мэйлом.
	_, err := s.userRepo.GetUserByUsernameAndPassword(ctx, newUser.Username, newUser.Email)
	if err == nil {
		return 0, errors.New("username or email already used")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	// проверяем, что пароль совпадает с подтверждением пароля
	if newUser.Password != newUser.PasswordConfirmation {
		return 0, errors.New("password and confirmation does not match")
	}

	// Хеширование пароля - никогда не храните пароль в незашифрованном виде.
	hashedPassword, err := GeneratePassword(s.passConfig, newUser.Password)
	if err != nil {
		return 0, err
	}

	// укладываем хэш пароля вместо изначальноего представления
	newUser.Password = hashedPassword

	// Передаем данные в слой репозитория для сохранения пользователя.
	userID, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return 0, err
	}

	// возвращаем данные в слой хэндлера
	return userID, nil
}

func (s *UserService) UpdateUser(ctx context.Context, updatedUser *models.UserDTO) (*models.UserDTO, error) {
	// Проверка наличия пользователя.
	existingUser, err := s.userRepo.GetUserByID(ctx, updatedUser.ID)
	if err != nil {
		return nil, err
	}

	// Обновление информации о пользователе.
	existingUser.Username = updatedUser.Username

	// Передаем данные в слой репозитория
	err = s.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	// возвращаем данные в слой хэндлера
	return updatedUser, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, request *models.UpdateUserPasswordDTO) error {
	// Проверка наличия пользователя.
	existingUser, err := s.userRepo.GetUserByID(ctx, request.ID)
	if err != nil {
		return err
	}

	// проверяем, совпадают ли новый пароль и подтверждение пароля
	if request.Password != request.PasswordConfirmation {
		return errors.New("password and confirmation does not match")
	}

	// Проверка старого пароля.
	passMatch, err := ComparePassword(request.OldPassword, existingUser.Password)
	if err != nil {
		return err
	}
	if !passMatch {
		return errors.New("incorrect old password")
	}

	// Хеширование пароля.
	hashedPassword, err := GeneratePassword(s.passConfig, request.Password)
	if err != nil {
		return err
	}

	// Обновление пароля в базе данных
	err = s.userRepo.UpdatePassword(ctx, request.ID, hashedPassword)
	if err != nil {
		return err
	}

	// возвращение ответа
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	// Удаление пользователя.
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	// Получение пользователя по его идентификатору.
	var userResponse = new(models.UserDTO)
	storedUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return userResponse, err
	}

	// запись данных из DAO - data access object через которую мы работаем с базой данных
	// в DTO - data transfer object, который мы возвращаем пользователю
	userResponse.ID = storedUser.ID
	userResponse.Username = storedUser.Username
	userResponse.Email = storedUser.Email

	// возврат данных пользователю
	return userResponse, nil
}

// GeneratePassword создает пароль на основе библиотеки golang.org/x/crypto/argon2
func GeneratePassword(c *config.PasswordConfig, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.Memory, c.Time, c.Threads, b64Salt, b64Hash)
	return full, nil
}

// ComparePassword сравниваем пароль и переданный хэш пароля на основе библиотеки golang.org/x/crypto/argon2
func ComparePassword(password, hash string) (bool, error) {

	parts := strings.Split(hash, "$")

	c := &config.PasswordConfig{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.KeyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}

func (s *UserService) Login(ctx context.Context, login *models.UserLoginDTO) (*models.UserTokens, error) {
	// Проверка наличия пользователя.
	existingUser, err := s.userRepo.GetUserByUsernameAndPassword(ctx, login.Username, login.Password)
	if err != nil {
		return nil, err
	}

	// Проверка пароля.
	passMatch, err := ComparePassword(login.Password, existingUser.Password)
	if err != nil {
		return nil, err
	}
	if !passMatch {
		return nil, errors.New("incorrect password")
	}

	// Генерируем токены и возвращаем
	accessToken, err := s.jwtUtil.GenerateAccessToken(existingUser.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtUtil.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		return nil, err
	}

	return &models.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Refresh(ctx context.Context, refresh string) (*models.UserTokens, error) {
	userID, err := s.jwtUtil.VerifyToken(refresh)
	if err != nil {
		return nil, err
	}

	// Генерируем токены и возвращаем
	accessToken, err := s.jwtUtil.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtUtil.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) VerifyToken(ctx context.Context, accessToken string) (int, error) {
	userID, err := s.jwtUtil.VerifyToken(accessToken)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
