package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO users (full_name, email, password_hash, university)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, email, password_hash, university, created_at, updated_at
	`

	created := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.University,
	).Scan(
		&created.ID,
		&created.FullName,
		&created.Email,
		&created.PasswordHash,
		&created.University,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT id, full_name, email, password_hash, university, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.University,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}
