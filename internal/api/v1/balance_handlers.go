package v1

import (
	"context"
	"finance/domain/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Balance response types
type BalanceResponse struct {
	AccountID        string           `json:"account_id"`
	CurrentBalance   string           `json:"current_balance"`
	PendingBalance   string           `json:"pending_balance"`
	AvailableBalance string           `json:"available_balance"`
	LastCalculated   string           `json:"last_calculated"`
	Account          *AccountResponse `json:"account,omitempty"`
}

type BalanceSummaryResponse struct {
	TotalAssets      string `json:"total_assets"`
	TotalLiabilities string `json:"total_liabilities"`
	NetWorth         string `json:"net_worth"`
	LastCalculated   string `json:"last_calculated"`
}

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/balance_uc.go . BalanceUseCase
type BalanceUseCase interface {
	GetBalanceByAccountID(ctx context.Context, accountID string) (entities.Balance, error)
	GetAllBalances(ctx context.Context) ([]entities.Balance, error)
	RefreshAccountBalance(ctx context.Context, accountID string) error
	RefreshAllBalances(ctx context.Context) error
	GetBalanceSummary(ctx context.Context) (entities.BalanceSummary, error)
}

// Balance handlers

// GetBalanceByAccountID retrieves the balance for a specific account
//
//	@Summary		Get balance by account ID
//	@Description	Retrieve the balance information for a specific account
//	@Tags			balances
//	@Accept			json
//	@Produce		json
//	@Param			accountId	path		string			true	"Account ID"
//	@Success		200			{object}	BalanceResponse	"Balance retrieved successfully"
//	@Failure		400			{object}	ErrorResponseBody	"Bad request"
//	@Failure		404			{object}	ErrorResponseBody	"Balance not found"
//	@Router			/balances/{accountId} [get]
func (h *ApiHandlers) GetBalanceByAccountID(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountId")
	if accountID == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("accountId"))
		return
	}

	balance, err := h.BalanceUseCase.GetBalanceByAccountID(r.Context(), accountID)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	response := BalanceResponse{
		AccountID:        balance.AccountID,
		CurrentBalance:   balance.CurrentBalance.String(),
		PendingBalance:   balance.PendingBalance.String(),
		AvailableBalance: balance.AvailableBalance.String(),
		LastCalculated:   balance.LastCalculated.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add account information if available
	if balance.Account != nil {
		response.Account = &AccountResponse{
			ID:          balance.Account.ID,
			Name:        balance.Account.Name,
			Type:        balance.Account.Type,
			Description: balance.Account.Description,
			CreatedAt:   balance.Account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   balance.Account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	render.JSON(w, r, response)
}

// GetAllBalances retrieves all account balances
//
//	@Summary		Get all balances
//	@Description	Retrieve a list of all account balances
//	@Tags			balances
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		BalanceResponse		"Balances retrieved successfully"
//	@Failure		500	{object}	ErrorResponseBody	"Internal server error"
//	@Router			/balances [get]
func (h *ApiHandlers) GetAllBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := h.BalanceUseCase.GetAllBalances(r.Context())
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	responses := make([]BalanceResponse, len(balances))
	for i, balance := range balances {
		responses[i] = BalanceResponse{
			AccountID:        balance.AccountID,
			CurrentBalance:   balance.CurrentBalance.String(),
			PendingBalance:   balance.PendingBalance.String(),
			AvailableBalance: balance.AvailableBalance.String(),
			LastCalculated:   balance.LastCalculated.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Add account information if available
		if balance.Account != nil {
			responses[i].Account = &AccountResponse{
				ID:          balance.Account.ID,
				Name:        balance.Account.Name,
				Type:        balance.Account.Type,
				Description: balance.Account.Description,
				CreatedAt:   balance.Account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:   balance.Account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
	}

	render.JSON(w, r, responses)
}

// GetBalanceSummary retrieves the overall balance summary
//
//	@Summary		Get balance summary
//	@Description	Retrieve a summary of all account balances including total assets, liabilities, and net worth
//	@Tags			balances
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	BalanceSummaryResponse	"Balance summary retrieved successfully"
//	@Failure		500	{object}	ErrorResponseBody		"Internal server error"
//	@Router			/balances/summary [get]
func (h *ApiHandlers) GetBalanceSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.BalanceUseCase.GetBalanceSummary(r.Context())
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	response := BalanceSummaryResponse{
		TotalAssets:      summary.TotalAssets.String(),
		TotalLiabilities: summary.TotalLiabilities.String(),
		NetWorth:         summary.NetWorth.String(),
		LastCalculated:   summary.LastCalculated.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

// RefreshAccountBalance refreshes the balance for a specific account
//
//	@Summary		Refresh account balance
//	@Description	Recalculate and refresh the balance for a specific account
//	@Tags			balances
//	@Accept			json
//	@Produce		json
//	@Param			accountId	path	string	true	"Account ID"
//	@Success		204			"Balance refreshed successfully"
//	@Failure		400			{object}	ErrorResponseBody	"Bad request"
//	@Failure		404			{object}	ErrorResponseBody	"Account not found"
//	@Router			/balances/{accountId}/refresh [post]
func (h *ApiHandlers) RefreshAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountId")
	if accountID == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("accountId"))
		return
	}

	err := h.BalanceUseCase.RefreshAccountBalance(r.Context(), accountID)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
