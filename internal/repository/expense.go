package repository

import (
	"context"
	"expense_tracker/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ExpenseRepository struct {
	db *Database
}

func NewExpenseRepository(db *Database) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

func (r *ExpenseRepository) CreateExpense(ctx context.Context, expense model.Expense) (int, error) {
	var id int
	q := `INSERT INTO expenses (user_id, amount, category, description, date)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.Pool.QueryRow(ctx, q, expense.UserID, expense.Amount, expense.Category, expense.Description, expense.Date).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository/expense: can't create expense: %w", err)
	}

	return id, nil
}

func (r *ExpenseRepository) GetExpenseByID(ctx context.Context, id int, userID int) (*model.Expense, error) {
	expense := &model.Expense{}
	q := `SELECT id, user_id, amount, category, description, date FROM expenses WHERE id = $1 and user_id = $2`
	err := r.db.Pool.QueryRow(ctx, q, id, userID).Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.Category, &expense.Description, &expense.Date)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("repository/expense: no such expense: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("repository/expense: can't found expense: %w", err)
	}
	return expense, nil
}

func (r *ExpenseRepository) IsExists(ctx context.Context, id int, userID int) (bool, error) {
	var count int
	q := `SELECT COUNT(*) FROM expenses WHERE id = $1 AND user_id = $2`
	err := r.db.Pool.QueryRow(ctx, q, id, userID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("repository/expense: can't check existanse of the expense: %w", err)
	}
	return count > 0, nil
}

func (r *ExpenseRepository) UpdateExpense(ctx context.Context, id int, userID int, input *model.UpdateExpenseInput) (*model.Expense, error) {
	updated := &model.Expense{}
	q := `UPDATE expenses SET amount = COALESCE($1, amount), category = COALESCE($2, category), 
	description = COALESCE($3, description), date = COALESCE($4, date) WHERE id = $5 and 
	user_id = $6 RETURNING id, user_id, amount, category, description, date`

	err := r.db.Pool.QueryRow(ctx, q, input.Amount, input.Category,
		input.Description, input.Date, id, userID).Scan(&updated.ID, &updated.UserID,
		&updated.Amount, &updated.Category, &updated.Description, &updated.Date)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("repository/expense: no such expense to update: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("repository/expense: can't update expense: %w", err)
	}

	return updated, nil
}

func (r *ExpenseRepository) DeleteExpense(ctx context.Context, id int, userID int) error {
	q := `DELETE FROM expanses WHERE id = $1 AND user_id = $2`
	result, err := r.db.Pool.Exec(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("repository/expense: can`t delete expanse: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository/expense: expanse not found (id: %d, user_id: %d)", id, userID)
	}
	return nil
}
