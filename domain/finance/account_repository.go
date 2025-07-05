package finance

import (
	"context"
	"finance/domain/entities"
)

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/account_repository.go . AccountRepository
type AccountRepository interface {
	CreateAccount(ctx context.Context, account entities.Account) (entities.Account, error)
	GetAccountByID(ctx context.Context, id string) (entities.Account, error)
	GetAllAccounts(ctx context.Context) ([]entities.Account, error)
	UpdateAccount(ctx context.Context, account entities.Account) (entities.Account, error)
	DeleteAccount(ctx context.Context, id string) error
	GetAccountWithBalance(ctx context.Context, id string) (entities.Account, error)
	GetAccountsWithBalances(ctx context.Context) ([]entities.Account, error)
}
