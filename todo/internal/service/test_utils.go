package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"todo/internal/service/mocks"
)

type Mocks struct {
	TodoRepository *mocks.MockTodoRepository
	RabbitProducer *mocks.MockRabbitProducer
}

func getMocks(ctrl *gomock.Controller) *Mocks {
	return &Mocks{
		TodoRepository: mocks.NewMockTodoRepository(ctrl),
		RabbitProducer: mocks.NewMockRabbitProducer(ctrl),
	}
}

func buildTestService(m *Mocks) *TodoService {
	return NewTodoService(
		m.TodoRepository,
		m.RabbitProducer,
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
