package v1

import (
	"context"
	"finance/domain/entities"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Use case interfaces
type AccountUseCase interface {
	CreateAccount(ctx context.Context, account entities.Account) (entities.Account, error)
	GetAccountByID(ctx context.Context, id string) (entities.Account, error)
	GetAllAccounts(ctx context.Context) ([]entities.Account, error)
	GetAccountsWithBalances(ctx context.Context) ([]entities.Account, error)
	UpdateAccount(ctx context.Context, account entities.Account) (entities.Account, error)
	DeleteAccount(ctx context.Context, id string) error
}

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error)
	GetCategoryByID(ctx context.Context, id string) (entities.Category, error)
	GetAllCategories(ctx context.Context) ([]entities.Category, error)
	GetCategoriesByType(ctx context.Context, categoryType entities.CategoryType) ([]entities.Category, error)
	UpdateCategory(ctx context.Context, category entities.Category) (entities.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (entities.Transaction, error)
	GetTransactionWithDetails(ctx context.Context, id string) (entities.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]entities.Transaction, error)
	GetTransactionsWithDetails(ctx context.Context, limit, offset int) ([]entities.Transaction, error)
	GetTransactionsByAccount(ctx context.Context, accountID string) ([]entities.Transaction, error)
	GetTransactionsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entities.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}

type BalanceUseCase interface {
	GetBalanceByAccountID(ctx context.Context, accountID string) (entities.Balance, error)
	GetAllBalances(ctx context.Context) ([]entities.Balance, error)
	RefreshAccountBalance(ctx context.Context, accountID string) error
	RefreshAllBalances(ctx context.Context) error
	GetBalanceSummary(ctx context.Context) (entities.BalanceSummary, error)
}

type ApiHandlers struct {
	ExampleUseCase     ExampleUseCase
	AccountUseCase     AccountUseCase
	CategoryUseCase    CategoryUseCase
	TransactionUseCase TransactionUseCase
	BalanceUseCase     BalanceUseCase
}

func (h *ApiHandlers) Routes(r chi.Router) {
	r.Get("/health", h.Health)
	r.Route("/api/v1", func(r chi.Router) {
		// Legacy example routes
		r.Route("/example", func(r chi.Router) {
			r.Post("/", h.CreateExample)
			r.Get("/{id}", h.GetExampleByID)
		})

		// Account routes
		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", h.CreateAccount)
			r.Get("/", h.GetAllAccounts)
			r.Get("/{id}", h.GetAccountByID)
			r.Put("/{id}", h.UpdateAccount)
			r.Delete("/{id}", h.DeleteAccount)
		})

		// Category routes
		r.Route("/categories", func(r chi.Router) {
			r.Post("/", h.CreateCategory)
			r.Get("/", h.GetAllCategories)
			r.Get("/{id}", h.GetCategoryByID)
			r.Put("/{id}", h.UpdateCategory)
			r.Delete("/{id}", h.DeleteCategory)
		})

		// Transaction routes
		r.Route("/transactions", func(r chi.Router) {
			r.Post("/", h.CreateTransaction)
			r.Get("/", h.GetAllTransactions)
			r.Get("/{id}", h.GetTransactionByID)
			r.Put("/{id}", h.UpdateTransaction)
			r.Delete("/{id}", h.DeleteTransaction)
		})

		// Balance routes
		r.Route("/balances", func(r chi.Router) {
			r.Get("/", h.GetAllBalances)
			r.Get("/summary", h.GetBalanceSummary)
			r.Get("/{accountId}", h.GetBalanceByAccountID)
			r.Post("/{accountId}/refresh", h.RefreshAccountBalance)
		})
	})
}

func (h *ApiHandlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type ErrorResponseBody struct {
	Error string `json:"error"`
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	render.Status(r, code)
	render.JSON(w, r, ErrorResponseBody{
		Error: err.Error(),
	})
}

func unknownErrorResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
}

// Helper functions
func errMissingParameter(param string) error {
	return fmt.Errorf("missing required parameter: %s", param)
}

func errNotFound(resource string) error {
	return fmt.Errorf("%s not found", resource)
}

func errInvalidParameter(param, value string) error {
	return fmt.Errorf("invalid parameter %s: %s", param, value)
}
