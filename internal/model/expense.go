package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Expense struct {
	ID          int     `json:"id,omitempty"`
	UserID      int     `json:"user_id,omitempty"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        Date    `json:"date"`
}

type UpdateExpenseInput struct {
	Amount      *float64 `json:"amount,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Description *string  `json:"description,omitempty"`
	Date        *Date    `json:"date,omitempty"`
}

type Date struct {
	time.Time
}

// Реализация интерфейса driver.Valuer
func (d Date) Value() (driver.Value, error) {
	return d.Time.Format("2006-01-02"), nil // Форматируем дату для БД
}

// Реализация интерфейса sql.Scanner
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

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format("2006-01-02"))), nil
}
