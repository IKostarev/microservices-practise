package service

import (
	"errors"
	"gateway/internal/service/mocks"
	"gateway/pkg/jwtutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type Mocks struct {
	UsersServiceClient *mocks.MockUsersServiceClient
	TodoServiceClient  *mocks.MockTodoServiceClient
}

func getMocks(ctrl *gomock.Controller) *Mocks {
	return &Mocks{
		UsersServiceClient: mocks.NewMockUsersServiceClient(ctrl),
		TodoServiceClient:  mocks.NewMockTodoServiceClient(ctrl),
	}
}

func buildTestService(m *Mocks) *GatewayService {
	return NewGatewayService(
		&jwtutil.JWTUtil{
			SecretKey:       "key",
			AccessTokenExp:  100000,
			RefreshTokenExp: 100000,
		},
		m.TodoServiceClient,
		m.UsersServiceClient,
		nil,
	)
}

func requireEqualError(t *testing.T, actualErr, expectedErr error) {
	if expectedErr == nil {
		require.NoError(t, actualErr)
	} else {
		require.Error(t, actualErr)
		if actualErr != nil {
			require.True(t, errors.Is(actualErr, expectedErr), "expected error: %v, got: %v", expectedErr, actualErr)
		}
	}
}
