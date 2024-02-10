package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"testing"
	"todo/internal/models"
)

func BenchmarkTodoService_CreateToDo(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mocks := getMocks(ctrl)
	service := buildTestService(mocks)
	ctx := context.Background()
	newTodo := &models.CreateTodoDTO{}

	mocks.TodoRepository.EXPECT().CreateToDo(gomock.Any(), gomock.Any()).Return(uuid.New(), nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CreateToDo(ctx, newTodo)
		if err != nil {
			b.Fatal(err)
		}
	}
}
