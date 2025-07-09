package pg

import (
	"context"
	"database/sql"
	"finance/domain/entities"
	"finance/internal/repository/pg/gen"
	"math/big"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/guilhermebr/gox/monetary"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepository struct {
	queries *gen.Queries
	db      *pgxpool.Pool
}

func NewBalanceRepository(db *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{
		queries: gen.New(db),
		db:      db,
	}
}

func (r *BalanceRepository) GetBalanceByAccountID(ctx context.Context, accountID string) (entities.Balance, error) {
	uuid, err := uuid.FromString(accountID)
	if err != nil {
		return entities.Balance{}, err
	}

	result, err := r.queries.GetBalanceByAccountID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Balance{}, nil
		}
		return entities.Balance{}, err
	}

	// Get the account to retrieve the asset information
	account, err := r.queries.GetAccountByID(ctx, result.AccountID)
	if err != nil {
		return entities.Balance{}, err
	}

	asset, ok := monetary.FindAssetByName(account.Asset)
	if !ok {
		asset = monetary.USD // default fallback
	}

	currentBalance, err := monetary.NewMonetary(asset, big.NewInt(result.CurrentBalance))
	if err != nil {
		return entities.Balance{}, err
	}

	pendingBalance, err := monetary.NewMonetary(asset, big.NewInt(result.PendingBalance))
	if err != nil {
		return entities.Balance{}, err
	}

	availableBalance, err := monetary.NewMonetary(asset, big.NewInt(result.AvailableBalance))
	if err != nil {
		return entities.Balance{}, err
	}

	return entities.Balance{
		AccountID:        result.AccountID.String(),
		CurrentBalance:   *currentBalance,
		PendingBalance:   *pendingBalance,
		AvailableBalance: *availableBalance,
		LastCalculated:   result.LastCalculated,
	}, nil
}

func (r *BalanceRepository) GetAllBalances(ctx context.Context) ([]entities.Balance, error) {
	results, err := r.queries.GetAllBalances(ctx)
	if err != nil {
		return nil, err
	}

	balances := make([]entities.Balance, len(results))
	for i, result := range results {
		// Get the account to retrieve the asset information
		account, err := r.queries.GetAccountByID(ctx, result.AccountID)
		if err != nil {
			return nil, err
		}

		asset, ok := monetary.FindAssetByName(account.Asset)
		if !ok {
			asset = monetary.USD // default fallback
		}

		currentBalance, err := monetary.NewMonetary(asset, big.NewInt(result.CurrentBalance))
		if err != nil {
			return nil, err
		}

		pendingBalance, err := monetary.NewMonetary(asset, big.NewInt(result.PendingBalance))
		if err != nil {
			return nil, err
		}

		availableBalance, err := monetary.NewMonetary(asset, big.NewInt(result.AvailableBalance))
		if err != nil {
			return nil, err
		}

		balances[i] = entities.Balance{
			AccountID:        result.AccountID.String(),
			CurrentBalance:   *currentBalance,
			PendingBalance:   *pendingBalance,
			AvailableBalance: *availableBalance,
			LastCalculated:   result.LastCalculated,
		}
	}

	return balances, nil
}

func (r *BalanceRepository) RefreshAccountBalance(ctx context.Context, accountID string) error {
	uuid, err := uuid.FromString(accountID)
	if err != nil {
		return err
	}

	return r.queries.RefreshAccountBalance(ctx, uuid)
}

func (r *BalanceRepository) GetBalanceSummary(ctx context.Context) (entities.BalanceSummary, error) {
	result, err := r.queries.GetBalanceSummary(ctx)
	if err != nil {
		return entities.BalanceSummary{}, err
	}

	// Convert interface{} values to proper types
	totalAssets, _ := result.TotalAssets.(int64)
	totalLiabilities, _ := result.TotalLiabilities.(int64)
	netWorth, _ := result.NetWorth.(int64)
	lastCalculated, _ := result.LastCalculated.(time.Time)

	// For balance summary, we'll use USD as the default asset
	// In a real implementation, you might want to have a configurable base currency
	usd := monetary.USD

	totalAssetsMonetary, err := monetary.NewMonetary(usd, big.NewInt(totalAssets))
	if err != nil {
		return entities.BalanceSummary{}, err
	}

	totalLiabilitiesMonetary, err := monetary.NewMonetary(usd, big.NewInt(totalLiabilities))
	if err != nil {
		return entities.BalanceSummary{}, err
	}

	netWorthMonetary, err := monetary.NewMonetary(usd, big.NewInt(netWorth))
	if err != nil {
		return entities.BalanceSummary{}, err
	}

	return entities.BalanceSummary{
		TotalAssets:      *totalAssetsMonetary,
		TotalLiabilities: *totalLiabilitiesMonetary,
		NetWorth:         *netWorthMonetary,
		LastCalculated:   lastCalculated,
	}, nil
}
