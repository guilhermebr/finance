package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ApiHandlers struct {
	AccountUseCase     AccountUseCase
	CategoryUseCase    CategoryUseCase
	TransactionUseCase TransactionUseCase
	BalanceUseCase     BalanceUseCase
}

func (h *ApiHandlers) Routes(r chi.Router) {
	r.Get("/health", h.Health)
	r.Route("/api/v1", func(r chi.Router) {

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

// Health returns the health status of the service
//
//	@Summary		Health check
//	@Description	Check if the service is healthy and running
//	@Tags			health
//	@Accept			json
//	@Produce		plain
//	@Success		200	"Service is healthy"
//	@Router			/health [get]
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
