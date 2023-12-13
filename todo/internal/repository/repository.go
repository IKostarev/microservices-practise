package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"todo/internal/models"
)

type TodoRepository struct {
	conn *pgxpool.Pool
}

func NewTodoRepository(conn *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{conn: conn}
}

func (r *TodoRepository) CreateToDo(ctx context.Context, newTodo *models.TodoDAO) (uuid.UUID, error) {
	var resID uuid.UUID

	sql := `INSERT INTO 
    					todo (created_by, assignee, description, created_at, updated_at)
			VALUES ($1, $2, $3, now(), now()) RETURNING id`

	err := r.conn.QueryRow(ctx, sql,
		newTodo.CreatedBy, newTodo.Assignee, newTodo.Description, newTodo.CreatedAt, newTodo.UpdatedAt).Scan(&resID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("[CreateToDO repo] create - %w\n", err)
	}

	return resID, nil
}

func (r *TodoRepository) UpdateToDo(ctx context.Context, newTodo *models.TodoDAO) error {
	sql := `UPDATE todo SET created_by = $1, assignee = $2, description = $3, updated_at = now() WHERE id = $5`

	_, err := r.conn.Exec(ctx, sql,
		newTodo.CreatedBy, newTodo.Assignee, newTodo.Description, newTodo.UpdatedAt, newTodo.ID)

	if err != nil {
		return fmt.Errorf("[UpdateToDO repo] update -  %w\n", err)
	}

	return nil
}

func (r *TodoRepository) GetToDos(ctx context.Context, todoID uuid.UUID) ([]models.TodoDAO, error) {
	res := make([]models.TodoDAO, 0)

	queryBuilder := squirrel.
		Select("id", "created_by", "assignee", "description", "created_at", "updated_at").
		Where("todo")

	if todoID != uuid.Nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": todoID})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("[GetToDos repo] place holder formater - %w\n", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("[GetToDos repo] get todos - %w\n", err)
	}

	for rows.Next() {
		var dao models.TodoDAO

		err = rows.Scan(
			&dao.ID,
			&dao.CreatedBy,
			&dao.Assignee,
			&dao.Description,
			&dao.CreatedAt,
			&dao.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("[GetToDos repo] get todos - %w\n", err)
		}

		res = append(res, dao)
	}

	return res, nil
}

func (r *TodoRepository) GetToDo(ctx context.Context, todoID uuid.UUID) (*models.TodoDAO, error) {
	var dao models.TodoDAO

	sql := `SELECT
				id, created_by, assignee, description, created_at, updated_at
			FROM
				todo
			WHERE
			    id = $1`

	err := r.conn.QueryRow(ctx, sql, todoID).
		Scan(
			&dao.ID,
			&dao.CreatedBy,
			&dao.Assignee,
			&dao.Description,
			&dao.CreatedAt,
			&dao.UpdatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("[GetToDo repo] get todo -  %w\n", err)
	}

	return &dao, nil
}

func (r *TodoRepository) DeleteToDo(ctx context.Context, todoID uuid.UUID) error {
	sql := `DELETE FROM todo WHERE id = $1`
	if _, err := r.conn.Exec(ctx, sql, todoID); err != nil {
		return fmt.Errorf("[DeleteTodo repo] delete -  %w\n", err)
	}

	return nil
}
