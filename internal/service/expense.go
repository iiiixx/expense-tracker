package service

import (
	"context"
	"expense_tracker/internal/model"
	"expense_tracker/internal/repository"
	"fmt"
	"strings"
	"time"
)

// ExpenseService provides methods for expense management.
type ExpenseService struct {
	expenseRepository *repository.ExpenseRepository
}

// NewExpenseService create an instanse of ExpenseService.
func NewExpenseService(expenseRepository *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepository: expenseRepository,
	}
}

// CreateExpense create an expense.
func (s *ExpenseService) CreateExpense(ctx context.Context, userID int, expense model.Expense) (*model.Expense, error) {
	if expense.Amount <= 0 {
		return nil, fmt.Errorf("service/expense: amount must be positive")
	}

	if len(expense.Category) == 0 {
		return nil, fmt.Errorf("service/expense: category is required")
	}

	expense.UserID = userID

	id, err := s.expenseRepository.CreateExpense(ctx, expense)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't create expense: %w", err)
	}

	expense.ID = id
	return &expense, nil
}

// GetExpenses retrieves an expense by ID.
func (s *ExpenseService) GetExpense(ctx context.Context, userID, expenseID int) (*model.Expense, error) {
	expense, err := s.expenseRepository.GetExpenseByID(ctx, expenseID, userID)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't get expense: %w", err)
	}

	return expense, nil
}

// UpdateExpense updates atributes of expenses by ID.
func (s *ExpenseService) UpdateExpense(ctx context.Context, expenseID, userID int, input *model.UpdateExpenseInput) (*model.Expense, error) {
	exists, err := s.expenseRepository.IsExists(ctx, expenseID, userID)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't found this expense: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("service/expense: expense not found")
	}

	updated, err := s.expenseRepository.UpdateExpense(ctx, expenseID, userID, input)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't update expense: %w", err)
	}
	return updated, nil
}

// DeleteExpense delete expense by ID.
func (s *ExpenseService) DeleteExpense(ctx context.Context, expenseID, userID int) error {
	if err := s.expenseRepository.DeleteExpense(ctx, expenseID, userID); err != nil {
		return fmt.Errorf("service/expense: can't delete expense: %w", err)
	}
	return nil
}

// GetExpensesList retrieves all user's expenses by user ID.
func (s *ExpenseService) GetExpensesList(ctx context.Context, userID int) ([]model.Expense, error) {

	expenses, err := s.expenseRepository.GetExpensesList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't get list of expenses: %w", err)
	}

	return expenses, nil
}

// GetExpensesByPeriod retrieves expenses for a user within a specific date range.
func (s *ExpenseService) GetExpensesByPeriod(ctx context.Context, userID int, start, end time.Time) ([]model.Expense, error) {

	expenses, err := s.expenseRepository.GetExpensesByPeriod(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't get expenses by period: %w", err)
	}

	return expenses, nil
}

// GetExpensesByCategory retrieves expenses in a specific category.
func (s *ExpenseService) GetExpensesByCategory(ctx context.Context, userID int, category string) ([]model.Expense, error) {
	if strings.TrimSpace(category) == "" {
		return nil, fmt.Errorf("service/expense: category can't be empty")
	}

	expenses, err := s.expenseRepository.GetExpensesByCategory(ctx, userID, category)
	if err != nil {
		return nil, fmt.Errorf("service/expense: can't get expenses by category: %w", err)
	}
	return expenses, nil

}
