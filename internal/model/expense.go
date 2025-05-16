package model

import (
	"fmt"
	"strings"
	"time"
)

type Expense struct {
	ID          int        `json:"id,omitempty"`
	UserID      int        `json:"user_id,omitempty"`
	Amount      float64    `json:"amount"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Date        CustomDate `json:"date"`
}

type UpdateExpenseInput struct {
	Amount      *float64    `json:"amount,omitempty"`
	Category    *string     `json:"category,omitempty"`
	Description *string     `json:"description,omitempty"`
	Date        *CustomDate `json:"date,omitempty"`
}

type CustomDate time.Time

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	*cd = CustomDate(t)
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(cd).Format("2006-01-02"))), nil
}
