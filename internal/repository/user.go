package repository

import (
	"context"
	"expense_tracker/internal/model"
	"fmt"

	"github.com/jackc/pgx"
)

// UserRepository provides data access methods for user operations.
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	q := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := r.db.Pool.QueryRow(ctx, q, user.Username, user.Password).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("repository/user: can't create user: %w", err)
	}
	return nil
}

// GetUserByName retrieves a user by their username.
func (r *UserRepository) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	q := `SELECT id, username, password FROM users WHERE username = $1`
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, q, username).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("repository/user: can't get user by name: %w", err)
	}
	return &user, err
}

// GetUserById retrieves a user by their ID.
func (r *UserRepository) GetUserById(ctx context.Context, id int) (*model.User, error) {
	q := `SELECT id, username, password FROM users WHERE id = $1`
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("repository/user: can't get user by id: %w", err)
	}
	return &user, err
}

// DeleteUser removes a user from the database by ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	q := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, q, id)

	if err != nil {
		return fmt.Errorf("repository/user: can't delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository/user: user with id %d not found", id)
	}
	return nil
}

// UpdateUsername changes a user's username.
func (r *UserRepository) UpdateUsername(ctx context.Context, id int, input *model.UpdateUsernameInput) (*model.User, error) {
	if input.Username == "" {
		return nil, fmt.Errorf("repository/user: username cannot be empty")
	}

	updated := model.User{}
	q := `UPDATE users SET username = COALESCE($1, username) WHERE id = $2 RETURNING id, username`

	err := r.db.Pool.QueryRow(ctx, q, input.Username,
		id).Scan(&updated.ID, &updated.Username)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("repository/user: no such user to update: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("repository/user: can't update user: %w", err)
	}

	return &updated, nil
}

// IsExistsUser checks if a user with given ID exists.
func (r *UserRepository) IsExistsUser(ctx context.Context, id int) (bool, error) {
	var count int
	q := `SELECT COUNT(*) FROM users WHERE id = $1`
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("repository/user: can't check existanse of the user: %w", err)
	}
	return count > 0, nil
}
