package finance

import (
	"context"
	"finance/domain/entities"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/guilhermebr/gox/monetary"
)

type TransactionUseCase struct {
	transactionRepo TransactionRepository
	accountRepo     AccountRepository
	categoryRepo    CategoryRepository
	balanceRepo     BalanceRepository
}

func NewTransactionUseCase(transactionRepo TransactionRepository, accountRepo AccountRepository, categoryRepo CategoryRepository, balanceRepo BalanceRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		balanceRepo:     balanceRepo,
	}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
	// Validate input
	if err := uc.validateTransaction(transaction); err != nil {
		return entities.Transaction{}, err
	}

	// Verify account exists
	account, err := uc.accountRepo.GetAccountByID(ctx, transaction.AccountID)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get account: %w", err)
	}
	if account.ID == "" {
		return entities.Transaction{}, fmt.Errorf("account not found")
	}

	// Convert the transaction amount to the correct asset based on the account
	// The handlers pass a temporary USD amount, so we need to convert it
	transaction = uc.convertTransactionToAccountAsset(transaction, account)

	// Verify category exists
	category, err := uc.categoryRepo.GetCategoryByID(ctx, transaction.CategoryID)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get category: %w", err)
	}
	if category.ID == "" {
		return entities.Transaction{}, fmt.Errorf("category not found")
	}

	// Set default status if not provided
	if transaction.Status == "" {
		transaction.Status = entities.TransactionStatusCleared
	}

	// Set default date if not provided
	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	createdTransaction, err := uc.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// The database trigger will automatically update the balance
	// But we can also refresh it manually to ensure consistency
	_ = uc.balanceRepo.RefreshAccountBalance(ctx, transaction.AccountID)

	return createdTransaction, nil
}

func (uc *TransactionUseCase) GetTransactionByID(ctx context.Context, id string) (entities.Transaction, error) {
	if id == "" {
		return entities.Transaction{}, fmt.Errorf("transaction ID cannot be empty")
	}

	transaction, err := uc.transactionRepo.GetTransactionByID(ctx, id)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

func (uc *TransactionUseCase) GetTransactionWithDetails(ctx context.Context, id string) (entities.Transaction, error) {
	if id == "" {
		return entities.Transaction{}, fmt.Errorf("transaction ID cannot be empty")
	}

	transaction, err := uc.transactionRepo.GetTransactionWithDetails(ctx, id)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get transaction with details: %w", err)
	}

	return transaction, nil
}

func (uc *TransactionUseCase) GetAllTransactions(ctx context.Context) ([]entities.Transaction, error) {
	transactions, err := uc.transactionRepo.GetAllTransactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (uc *TransactionUseCase) GetTransactionsWithDetails(ctx context.Context, limit, offset int) ([]entities.Transaction, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	transactions, err := uc.transactionRepo.GetTransactionsWithDetails(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions with details: %w", err)
	}

	return transactions, nil
}

func (uc *TransactionUseCase) GetTransactionsByAccount(ctx context.Context, accountID string) ([]entities.Transaction, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	transactions, err := uc.transactionRepo.GetTransactionsByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by account: %w", err)
	}

	return transactions, nil
}

func (uc *TransactionUseCase) GetTransactionsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Transaction, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, fmt.Errorf("start date and end date cannot be empty")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	transactions, err := uc.transactionRepo.GetTransactionsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}

	return transactions, nil
}

func (uc *TransactionUseCase) UpdateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
	// Validate input
	if err := uc.validateTransaction(transaction); err != nil {
		return entities.Transaction{}, err
	}

	if transaction.ID == "" {
		return entities.Transaction{}, fmt.Errorf("transaction ID cannot be empty")
	}

	// Get the existing transaction to compare account changes
	existingTransaction, err := uc.transactionRepo.GetTransactionByID(ctx, transaction.ID)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get existing transaction: %w", err)
	}

	if existingTransaction.ID == "" {
		return entities.Transaction{}, fmt.Errorf("transaction not found")
	}

	// Verify account exists
	account, err := uc.accountRepo.GetAccountByID(ctx, transaction.AccountID)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get account: %w", err)
	}
	if account.ID == "" {
		return entities.Transaction{}, fmt.Errorf("account not found")
	}

	// Convert the transaction amount to the correct asset based on the account
	transaction = uc.convertTransactionToAccountAsset(transaction, account)

	// Verify category exists
	category, err := uc.categoryRepo.GetCategoryByID(ctx, transaction.CategoryID)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to get category: %w", err)
	}
	if category.ID == "" {
		return entities.Transaction{}, fmt.Errorf("category not found")
	}

	// Business logic for transaction amounts based on category type
	transaction = uc.adjustTransactionAmount(transaction, category)

	updatedTransaction, err := uc.transactionRepo.UpdateTransaction(ctx, transaction)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Refresh balances for affected accounts
	_ = uc.balanceRepo.RefreshAccountBalance(ctx, transaction.AccountID)
	if existingTransaction.AccountID != transaction.AccountID {
		_ = uc.balanceRepo.RefreshAccountBalance(ctx, existingTransaction.AccountID)
	}

	return updatedTransaction, nil
}

func (uc *TransactionUseCase) DeleteTransaction(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("transaction ID cannot be empty")
	}

	// Get the transaction to know which account balance to refresh
	transaction, err := uc.transactionRepo.GetTransactionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction.ID == "" {
		return fmt.Errorf("transaction not found")
	}

	err = uc.transactionRepo.DeleteTransaction(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	// Refresh account balance
	_ = uc.balanceRepo.RefreshAccountBalance(ctx, transaction.AccountID)

	return nil
}

func (uc *TransactionUseCase) validateTransaction(transaction entities.Transaction) error {
	if transaction.AccountID == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	if transaction.CategoryID == "" {
		return fmt.Errorf("category ID cannot be empty")
	}

	if transaction.Monetary.Amount.Sign() == 0 {
		return fmt.Errorf("transaction amount cannot be zero")
	}

	if strings.TrimSpace(transaction.Description) == "" {
		return fmt.Errorf("transaction description cannot be empty")
	}

	if transaction.Status != "" {
		validStatuses := []entities.TransactionStatus{
			entities.TransactionStatusPending,
			entities.TransactionStatusCleared,
			entities.TransactionStatusCancelled,
		}

		isValidStatus := false
		for _, validStatus := range validStatuses {
			if transaction.Status == validStatus {
				isValidStatus = true
				break
			}
		}

		if !isValidStatus {
			return fmt.Errorf("invalid transaction status: %s", transaction.Status)
		}
	}

	return nil
}

func (uc *TransactionUseCase) adjustTransactionAmount(transaction entities.Transaction, category entities.Category) entities.Transaction {
	// For expense categories, ensure amount is negative
	if category.Type == entities.CategoryTypeExpense && transaction.Monetary.Amount.Sign() > 0 {
		negatedAmount := new(big.Int).Neg(transaction.Monetary.Amount)
		convertedMonetary, err := monetary.NewMonetary(transaction.Monetary.Asset, negatedAmount)
		if err == nil {
			transaction.Monetary = *convertedMonetary
		}
	}

	// For income categories, ensure amount is positive
	if category.Type == entities.CategoryTypeIncome && transaction.Monetary.Amount.Sign() < 0 {
		positiveAmount := new(big.Int).Neg(transaction.Monetary.Amount)
		convertedMonetary, err := monetary.NewMonetary(transaction.Monetary.Asset, positiveAmount)
		if err == nil {
			transaction.Monetary = *convertedMonetary
		}
	}

	return transaction
}

func (uc *TransactionUseCase) convertTransactionToAccountAsset(transaction entities.Transaction, account entities.Account) entities.Transaction {
	// If the transaction monetary asset is already the same as the account asset, no conversion needed
	if transaction.Monetary.Asset.Asset == account.Asset.Asset {
		return transaction
	}

	// Convert the amount to the account's asset
	// For simplicity, we'll keep the same numeric value but change the asset
	// In a real-world scenario, you would need currency conversion rates
	convertedMonetary, err := monetary.NewMonetary(account.Asset, transaction.Monetary.Amount)
	if err != nil {
		// If conversion fails, return the original transaction
		return transaction
	}

	transaction.Monetary = *convertedMonetary
	return transaction
}
