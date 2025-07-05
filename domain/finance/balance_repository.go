package finance

import (
	"context"
	"finance/domain/entities"
)

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/balance_repository.go . BalanceRepository
type BalanceRepository interface {
	GetBalanceByAccountID(ctx context.Context, accountID string) (entities.Balance, error)
	GetAllBalances(ctx context.Context) ([]entities.Balance, error)
	RefreshAccountBalance(ctx context.Context, accountID string) error
	GetBalanceSummary(ctx context.Context) (entities.BalanceSummary, error)
}
