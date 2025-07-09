package entities

import (
	"time"

	"github.com/guilhermebr/gox/monetary"
)

// Balance represents the balance of an account
type Balance struct {
	AccountID        string            `json:"account_id" db:"account_id"`
	CurrentBalance   monetary.Monetary `json:"current_balance" db:"current_balance"`
	PendingBalance   monetary.Monetary `json:"pending_balance" db:"pending_balance"`
	AvailableBalance monetary.Monetary `json:"available_balance" db:"available_balance"`
	LastCalculated   time.Time         `json:"last_calculated" db:"last_calculated"`

	// Relationships
	Account *Account `json:"account,omitempty"`
}

// BalanceSummary represents a summary of all account balances
type BalanceSummary struct {
	TotalAssets      monetary.Monetary `json:"total_assets"`
	TotalLiabilities monetary.Monetary `json:"total_liabilities"`
	NetWorth         monetary.Monetary `json:"net_worth"`
	LastCalculated   time.Time         `json:"last_calculated"`
}
