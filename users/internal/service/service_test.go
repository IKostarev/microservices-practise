package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	appErrors "users/internal/app_errors"
	"users/internal/models"
)

func TestUserService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t) // Создаем контроллер для моков.
	defer ctrl.Finish()             // Гарантируем завершение контроллера по окончании теста.

	mocks := getMocks(ctrl)        // Получаем моки для зависимостей сервиса.
	svc := buildTestService(mocks) // Строим тестируемый сервис с помощью моков.

	// Подготовка данных для нового пользователя.
	newUser := &models.CreateUserDTO{
		Username:             "testuser",
		Email:                "testuser@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	// Данные существующего пользователя для имитации сценария "Пользователь уже существует".
	existingUser := &models.UserDAO{
		Username: "testuser",
		Email:    "testuser@example.com",
	}

	// Определение тестовых сценариев.
	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success", // Сценарий успешной регистрации.
			setup: func() {
				// Настройка ожидаемого поведения моков.
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), newUser.Username, newUser.Email).Return(nil, sql.ErrNoRows)
				mocks.PasswordUtils.EXPECT().GeneratePassword(gomock.Any(), newUser.Password).Return("hashedPassword", nil)
				mocks.UsersRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(1, nil)
				mocks.RabbitProducer.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "UserExistsError", // Сценарий ошибки "Пользователь уже существует".
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), newUser.Username, newUser.Email).Return(existingUser, nil)
			},
			wantErr:     true,
			expectedErr: appErrors.ErrUsernameOrEmailIsUsed,
		},
		// Добавьте сюда другие сценарии
	}

	// Выполнение тестовых сценариев.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := svc.RegisterUser(context.Background(), newUser)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	updatedUser := &models.UserDTO{
		ID:       1,
		Username: "updateduser",
		Email:    "updateduser@example.com",
	}

	existingUser := &models.UserDAO{
		ID:       1,
		Username: "existinguser",
		Email:    "existinguser@example.com",
	}

	errDb := errors.New("db error")

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), updatedUser.ID).Return(existingUser, nil)
				mocks.UsersRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "UserNotFoundError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), updatedUser.ID).Return(nil, sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "UpdateUserDBError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), updatedUser.ID).Return(existingUser, nil)
				mocks.UsersRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(errDb)
			},
			wantErr:     true,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := svc.UpdateUser(context.Background(), updatedUser)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	userID := 1
	updateReq := &models.UpdateUserPasswordDTO{
		ID:                   userID,
		OldPassword:          "oldPassword",
		Password:             "newPassword",
		PasswordConfirmation: "newPassword",
	}

	existingUser := &models.UserDAO{
		ID:       userID,
		Password: "hashedOldPassword",
	}

	errDb := errors.New("db error")

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mocks.PasswordUtils.EXPECT().ComparePassword(gomock.Any(), updateReq.OldPassword, existingUser.Password).Return(true, nil)
				mocks.PasswordUtils.EXPECT().GeneratePassword(gomock.Any(), updateReq.Password).Return("hashedNewPassword", nil)
				mocks.UsersRepository.EXPECT().UpdatePassword(gomock.Any(), userID, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "UserNotFoundError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "UpdatePasswordDBError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mocks.PasswordUtils.EXPECT().ComparePassword(gomock.Any(), updateReq.OldPassword, existingUser.Password).Return(true, nil)
				mocks.PasswordUtils.EXPECT().GeneratePassword(gomock.Any(), updateReq.Password).Return("hashedNewPassword", nil)
				mocks.UsersRepository.EXPECT().UpdatePassword(gomock.Any(), userID, gomock.Any()).Return(errDb)
			},
			wantErr:     true,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := svc.UpdatePassword(context.Background(), updateReq)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	userID := 1
	errDb := errors.New("db error")

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.UsersRepository.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "UserNotFoundError",
			setup: func() {
				mocks.UsersRepository.EXPECT().DeleteUser(gomock.Any(), userID).Return(sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: sql.ErrNoRows,
		},
		{
			name: "DeleteUserDBError",
			setup: func() {
				mocks.UsersRepository.EXPECT().DeleteUser(gomock.Any(), userID).Return(errDb)
			},
			wantErr:     true,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := svc.DeleteUser(context.Background(), userID)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	userID := 1
	existingUser := &models.UserDAO{
		ID:       userID,
		Username: "existinguser",
		Email:    "existinguser@example.com",
	}

	errDb := errors.New("db error")

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
			},
			wantErr: false,
		},
		{
			name: "UserNotFoundError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: appErrors.ErrNotFound,
		},
		{
			name: "GetUserByIDDBError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errDb)
			},
			wantErr:     true,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := svc.GetUserByID(context.Background(), userID)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_GetUserByUsernameOrEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	username, email := "testuser", "testuser@example.com"
	existingUser := &models.UserDAO{
		Username: username,
		Email:    email,
	}

	errDb := errors.New("db error")

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), username, email).Return(existingUser, nil)
			},
			wantErr: false,
		},
		{
			name: "UserNotFoundError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), username, email).Return(nil, sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: appErrors.ErrNotFound,
		},
		{
			name: "GetUserByUsernameOrEmailDBError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), username, email).Return(nil, errDb)
			},
			wantErr:     true,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := svc.GetUserByUsernameOrEmail(context.Background(), username, email)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	loginReq := &models.UserLoginDTO{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password",
	}

	existingUser := &models.UserDAO{
		ID:       1,
		Username: loginReq.Username,
		Email:    loginReq.Email,
		Password: "hashedPassword",
	}

	errDb := errors.New("db error")

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), loginReq.Username, loginReq.Email).Return(existingUser, nil)
				mocks.PasswordUtils.EXPECT().ComparePassword(gomock.Any(), loginReq.Password, existingUser.Password).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name: "WrongCredentialsError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), loginReq.Username, loginReq.Email).Return(nil, sql.ErrNoRows)
			},
			wantErr:     true,
			expectedErr: appErrors.ErrWrongCredentials,
		},
		{
			name: "LoginDBError",
			setup: func() {
				mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), loginReq.Username, loginReq.Email).Return(nil, errDb)
			},
			wantErr:     true,
			expectedErr: errDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := svc.Login(context.Background(), loginReq)
			requireEqualError(t, err, tt.expectedErr)
		})
	}
}
