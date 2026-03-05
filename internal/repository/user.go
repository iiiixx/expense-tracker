package repository

import (
	"context"
	"expense_tracker/internal/model"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

// UserRepository provides data access methods for user operations.
type UserRepository struct {
	db     *Database
	logger *logrus.Logger
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *Database, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	log := r.logger.WithFields(logrus.Fields{
		"method":   "CreateUser",
		"username": user.Username,
	})
	log.Debug("creating new user")

	q := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err := r.db.Pool.QueryRow(ctx, q, user.Username, user.Password).Scan(&user.ID)

	if err != nil {
		log.WithError(err).Error("failed to create user")
		return fmt.Errorf("repository/user: can't create user: %w", err)
	}

	log.WithField("user_id", user.ID).Info("user created successfully")
	return nil
}

// GetUserByName retrieves a user by their username.
func (r *UserRepository) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":   "GetUserByName",
		"username": username,
	})
	log.Debug("fetching user by username")

	q := `SELECT id, username, password FROM users WHERE username = $1`
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, q, username).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn("user not found")
		} else {
			log.WithError(err).Error("failed to get user by username")
		}
		return nil, fmt.Errorf("repository/user: can't get user by name: %w", err)
	}

	log.WithField("user_id", user.ID).Debug("user found")
	return &user, nil
}

// GetUserById retrieves a user by their ID.
func (r *UserRepository) GetUserById(ctx context.Context, id int) (*model.User, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":  "GetUserById",
		"user_id": id,
	})
	log.Debug("fetching user by ID")

	q := `SELECT id, username, password FROM users WHERE id = $1`
	user := model.User{}
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn("user not found")
		} else {
			log.WithError(err).Error("failed to get user by ID")
		}
		return nil, fmt.Errorf("repository/user: can't get user by id: %w", err)
	}

	log.Debug("user retrieved successfully")
	return &user, nil
}

// DeleteUser removes a user from the database by ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	log := r.logger.WithFields(logrus.Fields{
		"method":  "DeleteUser",
		"user_id": id,
	})
	log.Info("deleting user")

	q := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, q, id)

	if err != nil {
		log.WithError(err).Error("failed to delete user")
		return fmt.Errorf("repository/user: can't delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		log.Warn("user not found for deletion")
		return fmt.Errorf("repository/user: user with id %d not found", id)
	}

	log.Info("user deleted successfully")
	return nil
}

// UpdateUsername changes a user's username.
func (r *UserRepository) UpdateUsername(ctx context.Context, id int, input *model.UpdateUsernameInput) (*model.User, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":       "UpdateUsername",
		"user_id":      id,
		"new_username": input.Username,
	})
	log.Info("updating username")

	if input.Username == "" {
		log.Warn("attempted update with empty username")
		return nil, fmt.Errorf("repository/user: username cannot be empty")
	}

	updated := model.User{}
	q := `UPDATE users SET username = COALESCE($1, username) WHERE id = $2 RETURNING id, username`

	err := r.db.Pool.QueryRow(ctx, q, input.Username,
		id).Scan(&updated.ID, &updated.Username)

	if err == pgx.ErrNoRows {
		log.Warn("user not found for update")
		return nil, fmt.Errorf("repository/user: no such user to update: %w", err)
	}
	if err != nil {
		log.WithError(err).Error("failed to update username")
		return nil, fmt.Errorf("repository/user: can't update user: %w", err)
	}

	log.WithField("new_username", updated.Username).Info("username updated successfully")
	return &updated, nil
}

// IsExistsUser checks if a user with given ID exists.
func (r *UserRepository) IsExistsUser(ctx context.Context, id int) (bool, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":  "IsExistsUser",
		"user_id": id,
	})
	log.Debug("checking user existence")

	var count int
	q := `SELECT COUNT(*) FROM users WHERE id = $1`
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(&count)

	if err != nil {
		log.WithError(err).Error("failed to check user existence")
		return false, fmt.Errorf("repository/user: can't check existanse of the user: %w", err)
	}

	exists := count > 0
	log.WithField("exists", exists).Debug("user existense check completed")
	return exists, nil
}
