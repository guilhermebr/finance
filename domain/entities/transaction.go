package entities

import (
	"time"

	"github.com/guilhermebr/gox/monetary"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCleared   TransactionStatus = "cleared"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          string            `json:"id" db:"id"`
	AccountID   string            `json:"account_id" db:"account_id"`
	CategoryID  string            `json:"category_id" db:"category_id"`
	Monetary    monetary.Monetary `json:"monetary" db:"monetary"`
	Description string            `json:"description" db:"description"`
	Date        time.Time         `json:"date" db:"date"`
	Status      TransactionStatus `json:"status" db:"status"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`

	// Relationships (for JSON responses)
	Account  *Account  `json:"account,omitempty"`
	Category *Category `json:"category,omitempty"`
}
