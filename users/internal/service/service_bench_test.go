package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"
	"users/internal/models"
)

// BenchmarkRegisterUser тестирует производительность функции RegisterUser.
func BenchmarkRegisterUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	newUser := &models.CreateUserDTO{ /* fill with test data */ }

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(1, nil).AnyTimes()
	mocks.PasswordUtils.EXPECT().GeneratePassword(gomock.Any(), gomock.Any()).Return("hashedPassword", nil).AnyTimes()
	mocks.RabbitProducer.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.RegisterUser(ctx, newUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpdateUser тестирует производительность функции UpdateUser.
func BenchmarkUpdateUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	updatedUser := &models.UserDTO{ /* fill with test data */ }

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&models.UserDAO{}, nil).AnyTimes()
	mocks.UsersRepository.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.UpdateUser(ctx, updatedUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpdatePassword тестирует производительность функции UpdatePassword.
func BenchmarkUpdatePassword(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	updatePasswordRequest := &models.UpdateUserPasswordDTO{ /* fill with test data */ }

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&models.UserDAO{Password: "oldHashedPassword"}, nil).AnyTimes()
	mocks.PasswordUtils.EXPECT().ComparePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()
	mocks.PasswordUtils.EXPECT().GeneratePassword(gomock.Any(), gomock.Any()).Return("newHashedPassword", nil).AnyTimes()
	mocks.UsersRepository.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.UpdatePassword(ctx, updatePasswordRequest)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDeleteUser тестирует производительность функции DeleteUser.
func BenchmarkDeleteUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	userID := 1 // Пример идентификатора пользователя

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().DeleteUser(gomock.Any(), gomock.Eq(userID)).Return(nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.DeleteUser(ctx, userID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetUserByID тестирует производительность функции GetUserByID.
func BenchmarkGetUserByID(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	userID := 1 // Пример идентификатора пользователя

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Eq(userID)).Return(&models.UserDAO{}, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetUserByID(ctx, userID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetUserByUsernameOrEmail тестирует производительность функции GetUserByUsernameOrEmail.
func BenchmarkGetUserByUsernameOrEmail(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	username := "exampleUser"
	email := "user@example.com"

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetUserByUsernameOrEmail(ctx, username, email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLogin тестирует производительность функции Login.
func BenchmarkLogin(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// Инициализация моков для тестирования
	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	loginDTO := &models.UserLoginDTO{ /* fill with test data */ }

	// Устанавливаем ожидаемые вызовы моков
	mocks.UsersRepository.EXPECT().GetUserByUsernameOrEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.UserDAO{}, nil).AnyTimes()
	mocks.PasswordUtils.EXPECT().ComparePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.Login(ctx, loginDTO)
		if err != nil {
			b.Fatal(err)
		}
	}
}
