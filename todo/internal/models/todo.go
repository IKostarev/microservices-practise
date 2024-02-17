package models

import (
	"github.com/google/uuid"
	"time"
	"todo/pkg/grpc_stubs/todo"
)

type TodoDAO struct {
	ID          uuid.UUID `db:"id"`
	CreatedBy   int       `db:"created_by"`
	Assignee    int       `db:"assignee"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type TodoDTO struct {
	ID          uuid.UUID `json:"id,omitempty" example:"c0e708fa-a7df-4d9f-a1b8-a3bfe63c433c"`
	CreatedBy   int       `json:"created_by" example:"1"`
	Assignee    int       `json:"assignee" example:"2"`
	Description string    `json:"description" example:"todo description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewEmptyTodoDTO() *TodoDTO {
	return &TodoDTO{}
}

func (t *TodoDTO) ToGRPC() *todo.TodoDTO {
	return &todo.TodoDTO{
		Id:          int32(uuid.UUID.ID(t.ID)),
		CreatedBy:   int32(t.CreatedBy),
		Assignee:    int32(t.Assignee),
		Description: t.Description,
	}
}

func (t *TodoDTO) FromGRPC(in *todo.TodoDTO) *TodoDTO {
	t.ID = uuid.New()
	t.CreatedBy = int(in.CreatedBy)
	t.Assignee = int(in.Assignee)
	t.Description = in.Description
	return t
}

type CreateTodoDTO struct {
	ID          uuid.UUID `json:"id,omitempty" example:"c0e708fa-a7df-4d9f-a1b8-a3bfe63c433c"`
	CreatedBy   int       `json:"created_by" example:"1"`
	Assignee    int       `json:"assignee" example:"2"`
	Description string    `json:"description" example:"todo description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewEmptyCreateTodoDTO() *CreateTodoDTO {
	return &CreateTodoDTO{}
}

func (t *CreateTodoDTO) ToGRPC() *todo.CreateTodoDTO {
	return &todo.CreateTodoDTO{
		Id:          int32(uuid.UUID.ID(t.ID)),
		CreatedBy:   int32(t.CreatedBy),
		Assignee:    int32(t.Assignee),
		Description: t.Description,
	}
}

func (t *CreateTodoDTO) FromGRPC(in *todo.CreateTodoDTO) *CreateTodoDTO {
	t.ID = uuid.New()
	t.CreatedBy = int(in.CreatedBy)
	t.Assignee = int(in.Assignee)
	t.Description = in.Description
	return t
}

type UpdateTodoDTO struct {
	ID          uuid.UUID `json:"id,omitempty" example:"c0e708fa-a7df-4d9f-a1b8-a3bfe63c433c"`
	UpdatedBy   int       `json:"updated_by" example:"1"`
	Assignee    int       `json:"assignee" example:"2"`
	Description string    `json:"description" example:"todo description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewEmptyUpdateTodoDTO() *UpdateTodoDTO {
	return &UpdateTodoDTO{}
}

func (t *UpdateTodoDTO) ToGRPC() *todo.UpdateTodoDTO {
	return &todo.UpdateTodoDTO{
		Id:          int32(uuid.UUID.ID(t.ID)),
		Assignee:    int32(t.Assignee),
		Description: t.Description,
	}
}

func (t *UpdateTodoDTO) FromGRPC(in *todo.UpdateTodoDTO) *UpdateTodoDTO {
	t.ID = uuid.New()
	t.Assignee = int(in.Assignee)
	t.Description = in.Description
	return t
}

type GetTodosDTO struct {
	CreatedBy int       `json:"created_by" example:"1"`
	Assignee  int       `json:"assignee" example:"2"`
	DateFrom  time.Time `json:"date_from"`
	DateTo    time.Time `json:"date_to"`
}
