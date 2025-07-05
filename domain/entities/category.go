package entities

import (
	"time"
)

// CategoryType represents the type of category
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

// Category represents a transaction category
type Category struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Type        CategoryType `json:"type" db:"type"`
	Description string       `json:"description" db:"description"`
	Color       string       `json:"color" db:"color"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}
