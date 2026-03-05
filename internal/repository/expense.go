package repository

import (
	"context"
	"expense_tracker/internal/model"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// ExpenseRepository provides data access methods for expense operations.
type ExpenseRepository struct {
	db     *Database
	logger *logrus.Logger
}

// NewExpenseRepository creates a new instance of ExpenseRepository.
func NewExpenseRepository(db *Database, logger *logrus.Logger) *ExpenseRepository {
	return &ExpenseRepository{
		db:     db,
		logger: logger,
	}
}

// CreateExpense inserts a new expense record into the database.
func (r *ExpenseRepository) CreateExpense(ctx context.Context, expense model.Expense) (int, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":   "CreateExpense",
		"user_id":  expense.UserID,
		"amount":   expense.Amount,
		"category": expense.Category,
	})
	log.Debug("creating new expense")

	var id int
	q := `INSERT INTO expenses (user_id, amount, category, description, date)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.Pool.QueryRow(ctx, q, expense.UserID, expense.Amount, expense.Category, expense.Description, expense.Date).Scan(&id)

	if err != nil {
		log.WithError(err).Error("failed to create expense")
		return 0, fmt.Errorf("repository/expense: can't create expense: %w", err)
	}

	log.WithField("expense_id", id).Info("expense created successfully")
	return id, nil
}

// GetExpenseByID retrieves an expense by its ID and associated user ID.
func (r *ExpenseRepository) GetExpenseByID(ctx context.Context, id int, userID int) (*model.Expense, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":     "GetExpenseByID",
		"expense_id": id,
		"user_id":    userID,
	})
	log.Debug("fetching expense")

	expense := &model.Expense{}
	q := `SELECT id, user_id, amount, category, description, date FROM expenses WHERE id = $1 and user_id = $2`
	err := r.db.Pool.QueryRow(ctx, q, id, userID).Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.Category, &expense.Description, &expense.Date)

	if err == pgx.ErrNoRows {
		log.Warn("expense not found")
		return nil, fmt.Errorf("repository/expense: no such expense: %w", err)
	}
	if err != nil {
		log.WithError(err).Error("failed to get expense")
		return nil, fmt.Errorf("repository/expense: can't found expense: %w", err)
	}
	log.Debug("expense retrieved successfully")
	return expense, nil
}

// IsExists checks if an expense with given ID exists for specified user.
func (r *ExpenseRepository) IsExists(ctx context.Context, id int, userID int) (bool, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":     "IsExists",
		"expense_id": id,
		"user_id":    userID,
	})
	log.Debug("checking expense existence")

	var count int
	q := `SELECT COUNT(*) FROM expenses WHERE id = $1 AND user_id = $2`
	err := r.db.Pool.QueryRow(ctx, q, id, userID).Scan(&count)

	if err != nil {
		log.WithError(err).Error("failed to check expense existence")
		return false, fmt.Errorf("repository/expense: can't check existanse of the expense: %w", err)
	}

	exists := count > 0
	log.WithField("exists", exists).Debug("expense existence check comleted")
	return exists, nil
}

// UpdateExpense modifies an existing expense record.
func (r *ExpenseRepository) UpdateExpense(ctx context.Context, id int, userID int, input *model.UpdateExpenseInput) (*model.Expense, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":     "UpdateExpense",
		"expense_id": id,
		"user_id":    userID,
		"input":      fmt.Sprintf("%+v", input),
	})
	log.Debug("updating expense")

	updated := &model.Expense{}

	q := `UPDATE expenses SET amount = COALESCE($1, amount), category = COALESCE($2, category), 
	description = COALESCE($3, description), date = COALESCE($4, date) WHERE id = $5 and 
	user_id = $6 RETURNING id, user_id, amount, category, description, date`

	err := r.db.Pool.QueryRow(ctx, q, input.Amount, input.Category,
		input.Description, input.Date, id, userID).Scan(&updated.ID, &updated.UserID,
		&updated.Amount, &updated.Category, &updated.Description, &updated.Date)
	if err == pgx.ErrNoRows {
		log.Warn("expense not found for update")
		return nil, fmt.Errorf("repository/expense: no such expense to update: %w", err)
	}
	if err != nil {
		log.WithError(err).Error("failed to update expense")
		return nil, fmt.Errorf("repository/expense: can't update expense: %w", err)
	}

	log.Info("expense updated successfully")
	return updated, nil
}

// DeleteExpense removes an expense record by ID.
func (r *ExpenseRepository) DeleteExpense(ctx context.Context, id int, userID int) error {
	log := r.logger.WithFields(logrus.Fields{
		"method":     "DeleteExpense",
		"expense_id": id,
		"user_id":    userID,
	})
	log.Debug("deleting expense")

	q := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`
	result, err := r.db.Pool.Exec(ctx, q, id, userID)
	if err != nil {
		log.WithError(err).Error("failed to delete expense")
		return fmt.Errorf("repository/expense: can`t delete expanse: %w", err)
	}

	if result.RowsAffected() == 0 {
		log.Warn("expense not found for deletion")
		return fmt.Errorf("repository/expense: expanse not found (id: %d, user_id: %d)", id, userID)
	}

	log.Info("expense deleted successfully")
	return nil
}

// GetExpensesList retrieves all expenses for a specific user.
func (r *ExpenseRepository) GetExpensesList(ctx context.Context, userID int) ([]model.Expense, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":  "GetExpensesList",
		"user_id": userID,
	})
	log.Debug("fetching expense list")

	q := `SELECT id, user_id, amount, category, description, date FROM expenses
	WHERE user_id = $1`

	rows, err := r.db.Pool.Query(ctx, q, userID)
	if err != nil {
		log.WithError(err).Error("failed to get expenses list")
		return nil, fmt.Errorf("repository/expense: can't get list of expenses: %w", err)
	}
	defer rows.Close()

	expenses, err := scanExpenses(rows, r.logger)
	if err != nil {
		return nil, err
	}

	log.WithField("count", len(expenses)).Debug("expenses list retrieved")
	return expenses, nil
}

// GetExpensesByPeriod retrieves expenses for a user within a specific date range.
func (r *ExpenseRepository) GetExpensesByPeriod(ctx context.Context, userID int, start, end time.Time) ([]model.Expense, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":  "GetExpensesByPeriod",
		"user_id": userID,
		"start":   start.Format(time.RFC3339),
		"end":     end.Format(time.RFC3339),
	})
	log.Debug("fetching expense by period")

	q := `SELECT id, user_id, amount, category, description, date FROM expenses
	WHERE user_id = $1 and date BETWEEN $2 AND $3 ORDER BY date`

	rows, err := r.db.Pool.Query(ctx, q, userID, start, end)
	if err != nil {
		log.WithError(err).Error("faied to get expenses by period")
		return nil, fmt.Errorf("repository/expense: can't get expenses by period: %w", err)
	}
	defer rows.Close()

	expenses, err := scanExpenses(rows, r.logger)
	if err != nil {
		return nil, err
	}

	log.WithField("count", len(expenses)).Debug("expenses by period retrieved")
	return expenses, nil
}

// GetExpensesByCategory retrieves expenses for a user in a specific category.
func (r *ExpenseRepository) GetExpensesByCategory(ctx context.Context, userID int, category string) ([]model.Expense, error) {
	log := r.logger.WithFields(logrus.Fields{
		"method":   "GetExpensesByCategory",
		"user_id":  userID,
		"category": category,
	})
	log.Debug("fetching expense by category")

	q := `SELECT id, user_id, amount, category, description, date FROM expenses
	WHERE user_id = $1 and category = $2 ORDER BY date`

	rows, err := r.db.Pool.Query(ctx, q, userID, category)
	if err != nil {
		log.WithError(err).Error("faied to get expenses by category")
		return nil, fmt.Errorf("repository/expense: can't get expenses by category: %w", err)
	}
	defer rows.Close()

	expenses, err := scanExpenses(rows, r.logger)
	if err != nil {
		return nil, err
	}

	log.WithField("count", len(expenses)).Debug("expenses by category retrieved")
	return expenses, nil
}

// scanExpenses is a helper function to scan multiple expense rows from a query result.
func scanExpenses(rows pgx.Rows, logger *logrus.Logger) ([]model.Expense, error) {
	log := logger.WithField("method", "scanExpenses")
	var expenses []model.Expense

	rowCount := 0
	for rows.Next() {
		rowCount++
		var e model.Expense
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.Amount,
			&e.Category,
			&e.Description,
			&e.Date,
		)
		if err != nil {
			log.WithError(err).Error("failed to scan expense row")
			return nil, fmt.Errorf("repository/expense: can't scan expense row: %w", err)
		}
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		log.WithError(err).Error("row iteration error")
		return nil, fmt.Errorf("repository/expense: rows ineration error: %w", err)
	}

	log.WithField("rows_scanned", rowCount).Debug("expenses scanning completed")
	return expenses, nil
}
