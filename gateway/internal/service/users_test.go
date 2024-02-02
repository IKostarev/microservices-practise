package service

import (
	"context"
	"gateway/internal/models"
	"gateway/pkg/ctxutil"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	userID := 1
	userToCreate := models.CreateUserDTO{}
	tests := []struct {
		name     string
		userData *models.CreateUserDTO
		setup    func(m *Mocks)
		wantErr  error
	}{
		{"success",
			&userToCreate,
			func(m *Mocks) {
				m.UsersServiceClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(userID, nil)
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svcMocks := getMocks(gomock.NewController(t))
			svc := buildTestService(svcMocks)

			test.setup(svcMocks)
			_, err := svc.RegisterUser(context.Background(), test.userData)

			requireEqualError(t, err, test.wantErr)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	updatedUser := &models.UserDTO{ID: 1}
	tests := []struct {
		name     string
		userData *models.UserDTO
		setup    func(m *Mocks)
		wantErr  error
	}{
		{
			"success",
			updatedUser,
			func(m *Mocks) {
				m.UsersServiceClient.EXPECT().UpdateUser(gomock.Any(), updatedUser).Return(nil)
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svcMocks := getMocks(gomock.NewController(t))
			svc := buildTestService(svcMocks)

			ctx := ctxutil.SetUserIDToContext(context.Background(), 1)

			test.setup(svcMocks)
			_, err := svc.UpdateUser(ctx, test.userData)

			requireEqualError(t, err, test.wantErr)
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	request := &models.UpdateUserPasswordDTO{ID: 1}
	tests := []struct {
		name    string
		request *models.UpdateUserPasswordDTO
		setup   func(m *Mocks)
		wantErr error
	}{
		{
			"success",
			request,
			func(m *Mocks) {
				m.UsersServiceClient.EXPECT().UpdatePassword(gomock.Any(), request).Return(nil)
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svcMocks := getMocks(gomock.NewController(t))
			svc := buildTestService(svcMocks)

			ctx := ctxutil.SetUserIDToContext(context.Background(), 1)

			test.setup(svcMocks)
			err := svc.UpdatePassword(ctx, test.request)

			requireEqualError(t, err, test.wantErr)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	userID := 1
	tests := []struct {
		name    string
		userID  int
		setup   func(m *Mocks)
		wantErr error
	}{
		{
			"success",
			userID,
			func(m *Mocks) {
				m.UsersServiceClient.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svcMocks := getMocks(gomock.NewController(t))
			svc := buildTestService(svcMocks)

			ctx := ctxutil.SetUserIDToContext(context.Background(), 1)

			test.setup(svcMocks)
			err := svc.DeleteUser(ctx, test.userID)

			requireEqualError(t, err, test.wantErr)
		})
	}
}
