package service

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	appErrors "users/app_errors"
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
		return 0, fmt.Errorf("[RegisterUser] get user: %w", appErrors.ErrUsernameOrEmailIsUsed)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	// проверяем, что пароль совпадает с подтверждением пароля
	if newUser.Password != newUser.PasswordConfirmation {
		return 0, fmt.Errorf("[RegisterUser] confirm pass: %w", appErrors.ErrPassAndConfirmationDoesNotMatch)
	}

	// Хеширование пароля - никогда не храните пароль в незашифрованном виде.
	hashedPassword, err := GeneratePassword(s.passConfig, newUser.Password)
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

	// возвращаем данные в слой хэндлера
	return userID, nil
}

func (s *UserService) UpdateUser(ctx context.Context, updatedUser *models.UserDTO) (*models.UserDTO, error) {
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
	passMatch, err := ComparePassword(request.OldPassword, existingUser.Password)
	if err != nil {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", err)
	}
	if !passMatch {
		return fmt.Errorf("[UpdatePassword] verify pass:%w", appErrors.ErrIncorrectOldPassword)
	}

	// Хеширование пароля.
	hashedPassword, err := GeneratePassword(s.passConfig, request.Password)
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
	// Удаление пользователя.
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("[DeleteUser] delete user:%w", err)
	}

	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	// Получение пользователя по его идентификатору.
	var userResponse = new(models.UserDTO)
	storedUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
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
		return nil, fmt.Errorf("[Login] get user:%w", err)
	}

	// Проверка пароля.
	passMatch, err := ComparePassword(login.Password, existingUser.Password)
	if err != nil {
		return nil, fmt.Errorf("[Login] verify pass:%w", err)
	}
	if !passMatch {
		return nil, fmt.Errorf("[Login] verify pass:%w", appErrors.ErrWrongPassword)
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

func (s *UserService) Refresh(ctx context.Context, refresh string) (*models.UserTokens, error) {
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

func (s *UserService) VerifyToken(ctx context.Context, accessToken string) (int, error) {
	userID, err := s.jwtUtil.VerifyToken(accessToken)
	if err != nil {
		return userID, fmt.Errorf("[VerifyToken] verify token:%w", err)
	}

	return userID, nil
}
