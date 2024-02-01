package service

import (
	"context"
	"gateway/internal/models"
	"github.com/golang/mock/gomock"
	"testing"
)

// BenchmarkRegisterUser тестирует производительность функции RegisterUser.
func BenchmarkRegisterUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков и настройка ожидаемых вызовов.
	mocks := getMocks(ctrl)
	mocks.UsersServiceClient.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		Return(1, nil).
		AnyTimes()

	// Сборка тестового сервиса с моками.
	service := buildTestService(mocks)

	// Тестовые данные для нового пользователя.
	newUser := &models.CreateUserDTO{ /* initialize with test data */ }

	// Сброс и запуск таймера бенчмарка.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.RegisterUser(context.Background(), newUser)
	}
}

// BenchmarkUpdateUser тестирует производительность функции UpdateUser.
func BenchmarkUpdateUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков.
	mocks := getMocks(ctrl)

	// Сборка тестового сервиса.
	service := buildTestService(mocks)

	// Тестовые данные для обновления пользователя.
	updatedUser := &models.UserDTO{ /* initialize with test data */ }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.UpdateUser(context.Background(), updatedUser)
	}
}

// BenchmarkUpdatePassword тестирует производительность функции UpdatePassword.
func BenchmarkUpdatePassword(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация и подготовка моков.
	mocks := getMocks(ctrl)

	// Сборка сервиса для тестирования.
	service := buildTestService(mocks)

	// Тестовые данные для запроса на обновление пароля.
	passwordUpdateRequest := &models.UpdateUserPasswordDTO{ /* initialize with test data */ }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.UpdatePassword(context.Background(), passwordUpdateRequest)
	}
}

// BenchmarkDeleteUser тестирует производительность функции DeleteUser.
func BenchmarkDeleteUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков и сборка сервиса.
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)

	// Пример идентификатора пользователя для удаления.
	userID := 1 // Example user ID

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.DeleteUser(context.Background(), userID)
	}
}

// BenchmarkGetUserByID тестирует производительность функции GetUserByID.
func BenchmarkGetUserByID(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Настройка ожидаемого вызова GetUserByID в моке.
	mocks := getMocks(ctrl)
	mocks.UsersServiceClient.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(&models.UserDTO{}, nil).
		AnyTimes()

	// Сборка сервиса с моками.
	service := buildTestService(mocks)

	// Идентификатор пользователя для теста.
	userID := 1 // Example user ID

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserByID(context.Background(), userID)
	}
}

// BenchmarkLogin тестирует производительность функции Login.
func BenchmarkLogin(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Настройка ожидаемого вызова UserLogin в моке.
	mocks := getMocks(ctrl)
	mocks.UsersServiceClient.EXPECT().
		UserLogin(gomock.Any(), gomock.Any()).
		Return(&models.UserDTO{}, nil).
		AnyTimes()

	// Сборка тестового сервиса.
	service := buildTestService(mocks)

	// Тестовые данные для входа в систему.
	loginRequest := &models.UserLoginDTO{ /* initialize with test data */ }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Login(context.Background(), loginRequest)
	}
}

// BenchmarkRefresh тестирует производительность функции Refresh.
func BenchmarkRefresh(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков и сборка сервиса.
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)

	// Тестовый refresh token.
	refreshToken := "example_refresh_token"
	accessToken := "example_refresh_token"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Refresh(context.Background(), refreshToken, accessToken)
	}
}
