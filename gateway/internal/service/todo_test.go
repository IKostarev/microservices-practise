package service

import (
	"context"
	appErrors "gateway/internal/app_errors"
	"gateway/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestTodoService_CreateToDo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	svc := buildTestService(mocks)

	newTodo := &models.CreateTodoDTO{
		ID:          uuid.New(),
		CreatedBy:   1,
		Assignee:    2,
		Description: "example desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	existingTodo := &models.TodoDAO{
		ID: uuid.MustParse("c0e708fa-a7df-4d9f-a1b8-a3bfe63c433c"),
	}

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success",
			setup: func() {
				mocks.TodoServiceClient.EXPECT().CreateToDo(gomock.Any(), newTodo.ID).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "TodoExistsError",
			setup: func() {
				mocks.TodoServiceClient.EXPECT().CreateToDo(gomock.Any(), newTodo).Return(existingTodo, nil)
			},
			wantErr:     true,
			expectedErr: appErrors.ErrUsernameOrEmailIsUsed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := svc.CreateToDo(context.Background(), newTodo)

			requireEqualError(t, err, tt.expectedErr)
		})
	}

}
