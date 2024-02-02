package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"users/internal/service/mocks"
)

type Mocks struct {
	UsersRepository *mocks.MockUserRepository
	RabbitProducer  *mocks.MockRabbitProducer
	PasswordUtils   *mocks.MockPasswordUtils
}

func getMocks(ctrl *gomock.Controller) *Mocks {
	return &Mocks{
		UsersRepository: mocks.NewMockUserRepository(ctrl),
		RabbitProducer:  mocks.NewMockRabbitProducer(ctrl),
		PasswordUtils:   mocks.NewMockPasswordUtils(ctrl),
	}
}

func buildTestService(m *Mocks) *UserService {
	return NewUserService(
		m.UsersRepository,
		m.RabbitProducer,
		m.PasswordUtils,
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
