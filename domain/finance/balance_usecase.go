package finance

import (
	"context"
	"finance/domain/entities"
	"fmt"
)

type BalanceUseCase struct {
	balanceRepo BalanceRepository
	accountRepo AccountRepository
}

func NewBalanceUseCase(balanceRepo BalanceRepository, accountRepo AccountRepository) *BalanceUseCase {
	return &BalanceUseCase{
		balanceRepo: balanceRepo,
		accountRepo: accountRepo,
	}
}

func (uc *BalanceUseCase) GetBalanceByAccountID(ctx context.Context, accountID string) (entities.Balance, error) {
	if accountID == "" {
		return entities.Balance{}, fmt.Errorf("account ID cannot be empty")
	}

	// Verify account exists
	account, err := uc.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return entities.Balance{}, fmt.Errorf("failed to get account: %w", err)
	}
	if account.ID == "" {
		return entities.Balance{}, fmt.Errorf("account not found")
	}

	balance, err := uc.balanceRepo.GetBalanceByAccountID(ctx, accountID)
	if err != nil {
		return entities.Balance{}, fmt.Errorf("failed to get balance: %w", err)
	}

	// If no balance record exists, create one
	if balance.AccountID == "" {
		err = uc.balanceRepo.RefreshAccountBalance(ctx, accountID)
		if err != nil {
			return entities.Balance{}, fmt.Errorf("failed to refresh account balance: %w", err)
		}

		balance, err = uc.balanceRepo.GetBalanceByAccountID(ctx, accountID)
		if err != nil {
			return entities.Balance{}, fmt.Errorf("failed to get balance after refresh: %w", err)
		}
	}

	// Add account information to the balance
	balance.Account = &account

	return balance, nil
}

func (uc *BalanceUseCase) GetAllBalances(ctx context.Context) ([]entities.Balance, error) {
	balances, err := uc.balanceRepo.GetAllBalances(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	// Enrich balances with account information
	for i := range balances {
		account, err := uc.accountRepo.GetAccountByID(ctx, balances[i].AccountID)
		if err != nil {
			// Skip accounts that can't be found
			continue
		}
		balances[i].Account = &account
	}

	return balances, nil
}

func (uc *BalanceUseCase) RefreshAccountBalance(ctx context.Context, accountID string) error {
	if accountID == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	// Verify account exists
	account, err := uc.accountRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account.ID == "" {
		return fmt.Errorf("account not found")
	}

	err = uc.balanceRepo.RefreshAccountBalance(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to refresh account balance: %w", err)
	}

	return nil
}

func (uc *BalanceUseCase) RefreshAllBalances(ctx context.Context) error {
	// Get all accounts
	accounts, err := uc.accountRepo.GetAllAccounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	// Refresh balance for each account
	for _, account := range accounts {
		err = uc.balanceRepo.RefreshAccountBalance(ctx, account.ID)
		if err != nil {
			// Log the error but continue with other accounts
			fmt.Printf("Failed to refresh balance for account %s: %v\n", account.ID, err)
		}
	}

	return nil
}

func (uc *BalanceUseCase) GetBalanceSummary(ctx context.Context) (entities.BalanceSummary, error) {
	summary, err := uc.balanceRepo.GetBalanceSummary(ctx)
	if err != nil {
		return entities.BalanceSummary{}, fmt.Errorf("failed to get balance summary: %w", err)
	}

	return summary, nil
}
