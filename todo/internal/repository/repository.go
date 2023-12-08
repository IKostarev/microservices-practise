package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"todo/internal/models"
)

type TodoRepository struct {
	conn *pgxpool.Pool
}

func NewTodoRepository(conn *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{conn: conn}
}

func (r *TodoRepository) CreateToDo(ctx context.Context, newTodo *models.TodoDAO) (uuid.UUID, error) {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	var resID uuid.UUID

	sql := `INSERT INTO 
    					todo (created_by, assignee, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.conn.QueryRow(context, sql,
		newTodo.CreatedBy, newTodo.Assignee, newTodo.Description, newTodo.CreatedAt, newTodo.UpdatedAt).Scan(&resID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("[CreateToDO repo] create - %w\n", err)
	}

	return resID, nil
}

func (r *TodoRepository) UpdateToDo(ctx context.Context, newTodo *models.TodoDAO) error {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	sql := `UPDATE todo SET created_by = $1, assignee = $2, description = $3, updated_at = $4 WHERE id = $5`

	_, err := r.conn.Exec(context, sql,
		newTodo.CreatedBy, newTodo.Assignee, newTodo.Description, newTodo.UpdatedAt, newTodo.ID)

	if err != nil {
		return fmt.Errorf("[UpdateToDO repo] update -  %w\n", err)
	}

	return nil
}

func (r *TodoRepository) GetToDos(ctx context.Context) ([]models.TodoDAO, error) {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	res := make([]models.TodoDAO, 0)

	sql := `SELECT id, created_by, assignee, description, created_at, updated_at FROM todo`

	rows, err := r.conn.Query(context, sql)
	if err != nil {
		return nil, fmt.Errorf("[GetToDos repo] get todos - %w\n", err)
	}

	var m models.TodoDAO

	for rows.Next() {
		err = rows.Scan(
			&m.ID,
			&m.CreatedBy,
			&m.Assignee,
			&m.Description,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("[GetToDos repo] get todos - %w\n", err)
		}

		res = append(res, m)
	}

	return res, nil
}

func (r *TodoRepository) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error) {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	var m models.TodoDAO

	sql := `SELECT
				id, created_by, assignee, description, created_at, updated_at
			FROM
				todo
			WHERE
			    id = $1`

	err := r.conn.QueryRow(context, sql, todoID).
		Scan(
			&m.ID,
			&m.CreatedBy,
			&m.Assignee,
			&m.Description,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo repo] get todo -  %w\n", err)
	}

	return &m, nil
}

func (r *TodoRepository) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	context, _ := context.WithTimeout(ctx, time.Second*3)

	sql := `DELETE FROM todo WHERE id = $1`
	if _, err := r.conn.Exec(context, sql, todoID); err != nil {
		return fmt.Errorf("[DeleteTodo repo] delete -  %w\n", err)
	}

	return nil
}
