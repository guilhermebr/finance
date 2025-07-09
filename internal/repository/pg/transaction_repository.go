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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	queries *gen.Queries
	db      *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		queries: gen.New(db),
		db:      db,
	}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
	accountID, err := uuid.FromString(transaction.AccountID)
	if err != nil {
		return entities.Transaction{}, err
	}

	categoryID, err := uuid.FromString(transaction.CategoryID)
	if err != nil {
		return entities.Transaction{}, err
	}

	date := pgtype.Date{
		Time:  transaction.Date,
		Valid: true,
	}

	// Convert monetary to int64 for storage
	amount := transaction.Monetary.Amount.Int64()

	result, err := r.queries.CreateTransaction(ctx, accountID, categoryID, amount, transaction.Description, date, string(transaction.Status))
	if err != nil {
		return entities.Transaction{}, err
	}

	// Get the account to retrieve the asset information
	account, err := r.queries.GetAccountByID(ctx, result.AccountID)
	if err != nil {
		return entities.Transaction{}, err
	}

	asset, ok := monetary.FindAssetByName(account.Asset)
	if !ok {
		asset = monetary.BRL // default fallback
	}

	// Convert back to monetary
	monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
	if err != nil {
		return entities.Transaction{}, err
	}

	return entities.Transaction{
		ID:          result.ID.String(),
		AccountID:   result.AccountID.String(),
		CategoryID:  result.CategoryID.String(),
		Monetary:    *monetaryAmount,
		Description: result.Description,
		Date:        result.Date.Time,
		Status:      entities.TransactionStatus(result.Status),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id string) (entities.Transaction, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return entities.Transaction{}, err
	}

	result, err := r.queries.GetTransactionByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Transaction{}, nil
		}
		return entities.Transaction{}, err
	}

	// Get the account to retrieve the asset information
	account, err := r.queries.GetAccountByID(ctx, result.AccountID)
	if err != nil {
		return entities.Transaction{}, err
	}

	asset, ok := monetary.FindAssetByName(account.Asset)
	if !ok {
		asset = monetary.BRL // default fallback
	}

	// Convert back to monetary
	monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
	if err != nil {
		return entities.Transaction{}, err
	}

	return entities.Transaction{
		ID:          result.ID.String(),
		AccountID:   result.AccountID.String(),
		CategoryID:  result.CategoryID.String(),
		Monetary:    *monetaryAmount,
		Description: result.Description,
		Date:        result.Date.Time,
		Status:      entities.TransactionStatus(result.Status),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *TransactionRepository) GetAllTransactions(ctx context.Context) ([]entities.Transaction, error) {
	results, err := r.queries.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	return r.convertTransactions(results), nil
}

func (r *TransactionRepository) GetTransactionsByAccount(ctx context.Context, accountID string) ([]entities.Transaction, error) {
	uuid, err := uuid.FromString(accountID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetTransactionsByAccount(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return r.convertTransactions(results), nil
}

func (r *TransactionRepository) GetTransactionsByCategory(ctx context.Context, categoryID string) ([]entities.Transaction, error) {
	uuid, err := uuid.FromString(categoryID)
	if err != nil {
		return nil, err
	}

	results, err := r.queries.GetTransactionsByCategory(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return r.convertTransactions(results), nil
}

func (r *TransactionRepository) GetTransactionsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Transaction, error) {
	startPgDate := pgtype.Date{Time: startDate, Valid: true}
	endPgDate := pgtype.Date{Time: endDate, Valid: true}

	results, err := r.queries.GetTransactionsByDateRange(ctx, startPgDate, endPgDate)
	if err != nil {
		return nil, err
	}

	return r.convertTransactions(results), nil
}

func (r *TransactionRepository) GetTransactionsByAccountAndDateRange(ctx context.Context, accountID string, startDate, endDate time.Time) ([]entities.Transaction, error) {
	uuid, err := uuid.FromString(accountID)
	if err != nil {
		return nil, err
	}

	startPgDate := pgtype.Date{Time: startDate, Valid: true}
	endPgDate := pgtype.Date{Time: endDate, Valid: true}

	results, err := r.queries.GetTransactionsByAccountAndDateRange(ctx, uuid, startPgDate, endPgDate)
	if err != nil {
		return nil, err
	}

	return r.convertTransactions(results), nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
	id, err := uuid.FromString(transaction.ID)
	if err != nil {
		return entities.Transaction{}, err
	}

	accountID, err := uuid.FromString(transaction.AccountID)
	if err != nil {
		return entities.Transaction{}, err
	}

	categoryID, err := uuid.FromString(transaction.CategoryID)
	if err != nil {
		return entities.Transaction{}, err
	}

	date := pgtype.Date{
		Time:  transaction.Date,
		Valid: true,
	}

	// Convert monetary to int64 for storage
	amount := transaction.Monetary.Amount.Int64()

	result, err := r.queries.UpdateTransaction(ctx, id, accountID, categoryID, amount, transaction.Description, date, string(transaction.Status))
	if err != nil {
		return entities.Transaction{}, err
	}

	// Get the account to retrieve the asset information
	account, err := r.queries.GetAccountByID(ctx, result.AccountID)
	if err != nil {
		return entities.Transaction{}, err
	}

	asset, ok := monetary.FindAssetByName(account.Asset)
	if !ok {
		asset = monetary.BRL // default fallback
	}

	// Convert back to monetary
	monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
	if err != nil {
		return entities.Transaction{}, err
	}

	return entities.Transaction{
		ID:          result.ID.String(),
		AccountID:   result.AccountID.String(),
		CategoryID:  result.CategoryID.String(),
		Monetary:    *monetaryAmount,
		Description: result.Description,
		Date:        result.Date.Time,
		Status:      entities.TransactionStatus(result.Status),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *TransactionRepository) UpdateTransactionStatus(ctx context.Context, id string, status entities.TransactionStatus) (entities.Transaction, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return entities.Transaction{}, err
	}

	result, err := r.queries.UpdateTransactionStatus(ctx, uuid, string(status))
	if err != nil {
		return entities.Transaction{}, err
	}

	// Get the account to retrieve the asset information
	account, err := r.queries.GetAccountByID(ctx, result.AccountID)
	if err != nil {
		return entities.Transaction{}, err
	}

	asset, ok := monetary.FindAssetByName(account.Asset)
	if !ok {
		asset = monetary.BRL // default fallback
	}

	// Convert back to monetary
	monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
	if err != nil {
		return entities.Transaction{}, err
	}

	return entities.Transaction{
		ID:          result.ID.String(),
		AccountID:   result.AccountID.String(),
		CategoryID:  result.CategoryID.String(),
		Monetary:    *monetaryAmount,
		Description: result.Description,
		Date:        result.Date.Time,
		Status:      entities.TransactionStatus(result.Status),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTransaction(ctx, uuid)
}

func (r *TransactionRepository) GetTransactionWithDetails(ctx context.Context, id string) (entities.Transaction, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return entities.Transaction{}, err
	}

	result, err := r.queries.GetTransactionWithDetails(ctx, uuid)
	if err != nil {
		return entities.Transaction{}, err
	}

	asset, ok := monetary.FindAssetByName(result.AccountAsset)
	if !ok {
		asset = monetary.BRL // default fallback
	}

	// Convert back to monetary
	monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
	if err != nil {
		return entities.Transaction{}, err
	}

	return entities.Transaction{
		ID:          result.ID.String(),
		AccountID:   result.AccountID.String(),
		CategoryID:  result.CategoryID.String(),
		Monetary:    *monetaryAmount,
		Description: result.Description,
		Date:        result.Date.Time,
		Status:      entities.TransactionStatus(result.Status),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		Account: &entities.Account{
			ID:   result.AccountID.String(),
			Name: result.AccountName,
			Type: entities.AccountType(result.AccountType),
		},
		Category: &entities.Category{
			ID:    result.CategoryID.String(),
			Name:  result.CategoryName,
			Type:  entities.CategoryType(result.CategoryType),
			Color: result.CategoryColor,
		},
	}, nil
}

func (r *TransactionRepository) GetTransactionsWithDetails(ctx context.Context, limit, offset int) ([]entities.Transaction, error) {
	results, err := r.queries.GetTransactionsWithDetails(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	transactions := make([]entities.Transaction, len(results))
	for i, result := range results {
		asset, ok := monetary.FindAssetByName(result.AccountAsset)
		if !ok {
			asset = monetary.BRL // default fallback
		}

		// Convert back to monetary
		monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
		if err != nil {
			continue // skip this transaction if we can't get the account
		}

		transactions[i] = entities.Transaction{
			ID:          result.ID.String(),
			AccountID:   result.AccountID.String(),
			CategoryID:  result.CategoryID.String(),
			Monetary:    *monetaryAmount,
			Description: result.Description,
			Date:        result.Date.Time,
			Status:      entities.TransactionStatus(result.Status),
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
			Account: &entities.Account{
				ID:   result.AccountID.String(),
				Name: result.AccountName,
				Type: entities.AccountType(result.AccountType),
			},
			Category: &entities.Category{
				ID:    result.CategoryID.String(),
				Name:  result.CategoryName,
				Type:  entities.CategoryType(result.CategoryType),
				Color: result.CategoryColor,
			},
		}
	}

	return transactions, nil
}

func (r *TransactionRepository) convertTransactions(results []gen.Transaction) []entities.Transaction {
	transactions := make([]entities.Transaction, len(results))
	for i, result := range results {
		// For convertTransactions, we need to fetch the account to get the asset
		// This is not ideal but necessary for the current implementation
		account, err := r.queries.GetAccountByID(context.Background(), result.AccountID)
		if err != nil {
			continue // skip this transaction if we can't get the account
		}

		asset, ok := monetary.FindAssetByName(account.Asset)
		if !ok {
			asset = monetary.BRL // default fallback
		}

		// Convert back to monetary
		monetaryAmount, err := monetary.NewMonetary(asset, big.NewInt(result.Amount))
		if err != nil {
			continue // skip this transaction if we can't get the account
		}

		transactions[i] = entities.Transaction{
			ID:          result.ID.String(),
			AccountID:   result.AccountID.String(),
			CategoryID:  result.CategoryID.String(),
			Monetary:    *monetaryAmount,
			Description: result.Description,
			Date:        result.Date.Time,
			Status:      entities.TransactionStatus(result.Status),
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return transactions
}
