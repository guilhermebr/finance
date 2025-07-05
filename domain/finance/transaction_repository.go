package finance

import (
	"context"
	"finance/domain/entities"
	"time"
)

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/transaction_repository.go . TransactionRepository
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (entities.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]entities.Transaction, error)
	GetTransactionsByAccount(ctx context.Context, accountID string) ([]entities.Transaction, error)
	GetTransactionsByCategory(ctx context.Context, categoryID string) ([]entities.Transaction, error)
	GetTransactionsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Transaction, error)
	GetTransactionsByAccountAndDateRange(ctx context.Context, accountID string, startDate, endDate time.Time) ([]entities.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id string, status entities.TransactionStatus) (entities.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
	GetTransactionWithDetails(ctx context.Context, id string) (entities.Transaction, error)
	GetTransactionsWithDetails(ctx context.Context, limit, offset int) ([]entities.Transaction, error)
}
