package pg

import (
	"context"
	"database/sql"
	"finance/domain/entities"
	"finance/internal/repository/pg/gen"
	"time"

	"github.com/gofrs/uuid/v5"
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

	return entities.Balance{
		AccountID:        result.AccountID.String(),
		CurrentBalance:   bigIntToFloat64(result.CurrentBalance),
		PendingBalance:   bigIntToFloat64(result.PendingBalance),
		AvailableBalance: bigIntToFloat64(result.AvailableBalance),
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
		balances[i] = entities.Balance{
			AccountID:        result.AccountID.String(),
			CurrentBalance:   bigIntToFloat64(result.CurrentBalance),
			PendingBalance:   bigIntToFloat64(result.PendingBalance),
			AvailableBalance: bigIntToFloat64(result.AvailableBalance),
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
	totalAssets, _ := result.TotalAssets.(float64)
	totalLiabilities, _ := result.TotalLiabilities.(float64)
	netWorth, _ := result.NetWorth.(float64)
	lastCalculated, _ := result.LastCalculated.(time.Time)

	return entities.BalanceSummary{
		TotalAssets:      totalAssets,
		TotalLiabilities: totalLiabilities,
		NetWorth:         netWorth,
		LastCalculated:   lastCalculated,
	}, nil
}
