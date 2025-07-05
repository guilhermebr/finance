package entities

import (
	"time"
)

// Balance represents the balance of an account
type Balance struct {
	AccountID        string    `json:"account_id" db:"account_id"`
	CurrentBalance   float64   `json:"current_balance" db:"current_balance"`
	PendingBalance   float64   `json:"pending_balance" db:"pending_balance"`
	AvailableBalance float64   `json:"available_balance" db:"available_balance"`
	LastCalculated   time.Time `json:"last_calculated" db:"last_calculated"`

	// Relationships
	Account *Account `json:"account,omitempty"`
}

// BalanceSummary represents a summary of all account balances
type BalanceSummary struct {
	TotalAssets      float64   `json:"total_assets"`
	TotalLiabilities float64   `json:"total_liabilities"`
	NetWorth         float64   `json:"net_worth"`
	LastCalculated   time.Time `json:"last_calculated"`
}
