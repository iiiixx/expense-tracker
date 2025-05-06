package repository

import (
	"context"
	"expense_tracker/internal/model"
	"fmt"

	"github.com/jackc/pgx"
)

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	q := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := r.db.Pool.QueryRow(ctx, q, user.Username, user.Password).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("repository: can't create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*model.User, error) {
	q := `SELECT id, username, password FROM users WHERE user_id = $1`
	user := &model.User{}
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("repository: can't get user: %w", err)
	}
	return user, err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	q := `DELETE FROM users WHERE id = &1`
	result, err := r.db.Pool.Exec(ctx, q, id)

	if err != nil {
		return fmt.Errorf("repository: can't delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository: user with id %d not found", id)
	}
	return nil
}

func (r *UserRepository) UpdateUsername(ctx context.Context, id int, newUserName string) (*model.User, error) {
	if newUserName == "" {
		return nil, fmt.Errorf("repository: username cannot be empty")
	}

	updated := &model.User{}
	q := `UPDATE users SET username = COALESCE($1, username) WHERE id = $2 RETURNING id, username`

	err := r.db.Pool.QueryRow(ctx, q, newUserName,
		id).Scan(&updated.ID, &updated.Username)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("repository: no such user to update: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("repository: can't update user: %w", err)
	}

	return updated, nil
}
