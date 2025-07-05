package pg

import (
	"context"
	"database/sql"
	"finance/domain/entities"
	"finance/internal/repository/pg/gen"
	"math/big"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	queries *gen.Queries
	db      *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		queries: gen.New(db),
		db:      db,
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account entities.Account) (entities.Account, error) {
	result, err := r.queries.CreateAccount(ctx, account.Name, string(account.Type), &account.Description)
	if err != nil {
		return entities.Account{}, err
	}

	return entities.Account{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.AccountType(result.Type),
		Description: ptrStringValue(result.Description),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, id string) (entities.Account, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return entities.Account{}, err
	}

	result, err := r.queries.GetAccountByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Account{}, nil
		}
		return entities.Account{}, err
	}

	return entities.Account{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.AccountType(result.Type),
		Description: ptrStringValue(result.Description),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *AccountRepository) GetAllAccounts(ctx context.Context) ([]entities.Account, error) {
	results, err := r.queries.GetAllAccounts(ctx)
	if err != nil {
		return nil, err
	}

	accounts := make([]entities.Account, len(results))
	for i, result := range results {
		accounts[i] = entities.Account{
			ID:          result.ID.String(),
			Name:        result.Name,
			Type:        entities.AccountType(result.Type),
			Description: ptrStringValue(result.Description),
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return accounts, nil
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, account entities.Account) (entities.Account, error) {
	uuid, err := uuid.FromString(account.ID)
	if err != nil {
		return entities.Account{}, err
	}

	result, err := r.queries.UpdateAccount(ctx, uuid, account.Name, string(account.Type), &account.Description)
	if err != nil {
		return entities.Account{}, err
	}

	return entities.Account{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.AccountType(result.Type),
		Description: ptrStringValue(result.Description),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *AccountRepository) DeleteAccount(ctx context.Context, id string) error {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteAccount(ctx, uuid)
}

func (r *AccountRepository) GetAccountWithBalance(ctx context.Context, id string) (entities.Account, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return entities.Account{}, err
	}

	result, err := r.queries.GetAccountWithBalance(ctx, uuid)
	if err != nil {
		return entities.Account{}, err
	}

	return entities.Account{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.AccountType(result.Type),
		Description: ptrStringValue(result.Description),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *AccountRepository) GetAccountsWithBalances(ctx context.Context) ([]entities.Account, error) {
	results, err := r.queries.GetAccountsWithBalances(ctx)
	if err != nil {
		return nil, err
	}

	accounts := make([]entities.Account, len(results))
	for i, result := range results {
		accounts[i] = entities.Account{
			ID:          result.ID.String(),
			Name:        result.Name,
			Type:        entities.AccountType(result.Type),
			Description: ptrStringValue(result.Description),
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return accounts, nil
}

// Helper functions
func ptrStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func bigIntToFloat64(bi big.Int) float64 {
	f, _ := bi.Float64()
	return f / 100.0 // Convert from cents to dollars
}

func float64ToBigInt(f float64) *big.Int {
	cents := int64(f * 100) // Convert to cents
	return big.NewInt(cents)
}
