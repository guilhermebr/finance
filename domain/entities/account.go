package entities

import (
	"time"

	"github.com/guilhermebr/gox/monetary"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCredit     AccountType = "credit"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeCash       AccountType = "cash"
)

// Account represents a financial account
type Account struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Type        AccountType    `json:"type" db:"type"`
	Asset       monetary.Asset `json:"asset" db:"asset"`
	Description string         `json:"description" db:"description"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}
