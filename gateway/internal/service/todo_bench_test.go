package service

import (
	"context"
	"gateway/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"testing"
)

func BenchmarkTodoService_CreateToDo(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	newTodo := &models.CreateTodoDTO{}

	mocks.TodoServiceClient.EXPECT().CreateToDo(gomock.Any(), gomock.Any()).Return(uuid.New(), nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CreateToDo(ctx, newTodo)
		if err != nil {
			b.Fatal(err)
		}
	}
}
