package finance

import (
	"context"
	"finance/domain/entities"
	"fmt"
	"log/slog"
	"strings"

	"github.com/guilhermebr/gox/monetary"
)

type AccountUseCase struct {
	accountRepo AccountRepository
	balanceRepo BalanceRepository
}

func NewAccountUseCase(accountRepo AccountRepository, balanceRepo BalanceRepository) *AccountUseCase {
	return &AccountUseCase{
		accountRepo: accountRepo,
		balanceRepo: balanceRepo,
	}
}

func (uc *AccountUseCase) CreateAccount(ctx context.Context, account entities.Account) (entities.Account, error) {
	// Validate input
	if err := uc.validateAccount(account); err != nil {
		return entities.Account{}, err
	}

	// Create the account
	createdAccount, err := uc.accountRepo.CreateAccount(ctx, account)
	if err != nil {
		return entities.Account{}, fmt.Errorf("failed to create account: %w", err)
	}

	// Initialize balance for the new account
	err = uc.balanceRepo.RefreshAccountBalance(ctx, createdAccount.ID)
	if err != nil {
		slog.Error("failed to refresh account balance")
		// Log the error but don't fail the account creation
		// The balance will be calculated when the first transaction is added
	}

	return createdAccount, nil
}

func (uc *AccountUseCase) GetAccountByID(ctx context.Context, id string) (entities.Account, error) {
	if id == "" {
		return entities.Account{}, fmt.Errorf("account ID cannot be empty")
	}

	account, err := uc.accountRepo.GetAccountByID(ctx, id)
	if err != nil {
		return entities.Account{}, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (uc *AccountUseCase) GetAllAccounts(ctx context.Context) ([]entities.Account, error) {
	accounts, err := uc.accountRepo.GetAllAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (uc *AccountUseCase) UpdateAccount(ctx context.Context, account entities.Account) (entities.Account, error) {
	// Validate input
	if err := uc.validateAccount(account); err != nil {
		return entities.Account{}, err
	}

	if account.ID == "" {
		return entities.Account{}, fmt.Errorf("account ID cannot be empty")
	}

	// Check if account exists
	existingAccount, err := uc.accountRepo.GetAccountByID(ctx, account.ID)
	if err != nil {
		return entities.Account{}, fmt.Errorf("failed to get existing account: %w", err)
	}

	if existingAccount.ID == "" {
		return entities.Account{}, fmt.Errorf("account not found")
	}

	updatedAccount, err := uc.accountRepo.UpdateAccount(ctx, account)
	if err != nil {
		return entities.Account{}, fmt.Errorf("failed to update account: %w", err)
	}

	return updatedAccount, nil
}

func (uc *AccountUseCase) DeleteAccount(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	// Check if account exists
	account, err := uc.accountRepo.GetAccountByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	if account.ID == "" {
		return fmt.Errorf("account not found")
	}

	// TODO: Check if account has transactions and handle accordingly
	// For now, we'll rely on the database cascade delete

	err = uc.accountRepo.DeleteAccount(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

func (uc *AccountUseCase) validateAccount(account entities.Account) error {
	if strings.TrimSpace(account.Name) == "" {
		return fmt.Errorf("account name cannot be empty")
	}

	if account.Type == "" {
		return fmt.Errorf("account type cannot be empty")
	}

	validTypes := []entities.AccountType{
		entities.AccountTypeChecking,
		entities.AccountTypeSavings,
		entities.AccountTypeCredit,
		entities.AccountTypeInvestment,
		entities.AccountTypeCash,
	}

	isValidType := false
	for _, validType := range validTypes {
		if account.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid account type: %s", account.Type)
	}

	// Validate asset
	if account.Asset.Asset == "" {
		return fmt.Errorf("account asset cannot be empty")
	}

	// Verify that the asset is valid by checking if it exists in the monetary package
	_, ok := monetary.FindAssetByName(account.Asset.Asset)
	if !ok {
		return fmt.Errorf("invalid asset: %s", account.Asset.Asset)
	}

	return nil
}
