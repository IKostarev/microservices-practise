package models

import (
	"github.com/google/uuid"
	"time"
)

type TodoDTO struct {
	ID          uuid.UUID `json:"id,omitempty" example:"c0e708fa-a7df-4d9f-a1b8-a3bfe63c433c"`
	CreatedBy   int       `json:"created_by" example:"1"`
	Assignee    int       `json:"assignee" example:"2"`
	Description string    `json:"description" example:"todo description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetTodosDTO struct {
	CreatedBy int       `json:"created_by" example:"1"`
	Assignee  int       `json:"assignee" example:"2"`
	DateFrom  time.Time `json:"date_from"`
	DateTo    time.Time `json:"date_to"`
}
