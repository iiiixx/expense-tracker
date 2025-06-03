package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Expense represents a financial expense record in the system.
type Expense struct {
	ID          int     `json:"id,omitempty"`
	UserID      int     `json:"user_id,omitempty"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        Date    `json:"date"`
}

// UpdateExpenseInput contains fields for updating an existing expense record. All fields are optional.
type UpdateExpenseInput struct {
	Amount      *float64 `json:"amount,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Description *string  `json:"description,omitempty"`
	Date        *Date    `json:"date,omitempty"`
}

// Custom date type that extends time.Time with specific serialization behavior.
type Date struct {
	time.Time
}

// Implements driver.Valuer for proper SQL storage (formats as YYYY-MM-DD)
func (d Date) Value() (driver.Value, error) {
	return d.Time.Format("2006-01-02"), nil // Форматируем дату для БД
}

// Implements sql.Scanner for reading from database
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("invalid type for Date: %T", value)
	}

	d.Time = t
	return nil
}

// UnmarshalJSON parses dates from JSON strings in YYYY-MM-DD format
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	d.Time = t
	return nil
}

// MarshalJSON formats dates as JSON strings in YYYY-MM-DD format
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format("2006-01-02"))), nil
}
