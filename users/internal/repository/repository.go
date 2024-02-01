package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"users/internal/models"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn}
}
func (r *UserRepository) CreateUser(ctx context.Context, user *models.CreateUserDTO) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.CreateUser")
	defer span.Finish()

	var userID int
	sql := `
        INSERT INTO 
			users (
			   username, 
			   password, 
			   email
			)
        VALUES 
			($1, $2, $3)
        RETURNING id
    `
	err := r.conn.QueryRow(ctx, sql, user.Username, user.Password, user.Email).
		Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.UserDAO) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UpdateUser")
	defer span.Finish()

	sql := `
        UPDATE 
            users
        SET 
			username = $2, 
			email = $3
        WHERE 
            id = $1
    `
	_, err := r.conn.Exec(ctx, sql, user.ID, user.Username, user.Email)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newPassword string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.UpdatePassword")
	defer span.Finish()

	sql := `
        UPDATE 
            users
        SET 
            password = $2
        WHERE 
            id = $1
    `
	_, err := r.conn.Exec(ctx, sql, userID, newPassword)
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.DeleteUser")
	defer span.Finish()

	sql := `
        DELETE FROM 
		    users
        WHERE 
            id = $1
    `
	_, err := r.conn.Exec(ctx, sql, userID)
	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*models.UserDAO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetUserByID")
	defer span.Finish()

	var user models.UserDAO
	sql := `
        SELECT 
            id, 
            username, 
            password, 
            email
        FROM 
            users
        WHERE 
            id = $1
    `
	err := r.conn.QueryRow(ctx, sql, userID).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsernameOrEmail(ctx context.Context, username, email string) (*models.UserDAO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetUserByUsernameOrEmail")
	defer span.Finish()

	var user models.UserDAO
	// создадим конструктор квери. определим, что и откуда забрать
	queryBuilder := sq.
		Select("id", "username", "password", "email").
		From("users")

	// убедимся, что username или email не являются пустыми строками перед добавлением их в запрос
	if username != "" && email == "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"username": username})
	}
	if email != "" && username == "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"email": email})
	}

	if email != "" && username != "" {
		queryBuilder = queryBuilder.Where(sq.Or{sq.Eq{"username": username}, sq.Eq{"email": email}})
	}

	// создадим квери и аргументы для нее, зададим формат плэйсхолдеров в виде доллара
	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	// выполним запрос с созданной квери и подготовленными аргументами
	err = r.conn.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.UserDAO, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.GetUserByUsername")
	defer span.Finish()

	var user models.UserDAO
	sql := `
        SELECT 
            id, 
            username, 
            password, 
            email
        FROM 
            users
        WHERE 
            username = $1
    `
	err := r.conn.QueryRow(ctx, sql, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
