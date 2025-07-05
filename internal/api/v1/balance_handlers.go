package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Balance response types
type BalanceResponse struct {
	AccountID        string           `json:"account_id"`
	CurrentBalance   float64          `json:"current_balance"`
	PendingBalance   float64          `json:"pending_balance"`
	AvailableBalance float64          `json:"available_balance"`
	LastCalculated   string           `json:"last_calculated"`
	Account          *AccountResponse `json:"account,omitempty"`
}

type BalanceSummaryResponse struct {
	TotalAssets      float64 `json:"total_assets"`
	TotalLiabilities float64 `json:"total_liabilities"`
	NetWorth         float64 `json:"net_worth"`
	LastCalculated   string  `json:"last_calculated"`
}

// Balance handlers
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
		CurrentBalance:   balance.CurrentBalance,
		PendingBalance:   balance.PendingBalance,
		AvailableBalance: balance.AvailableBalance,
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
			CurrentBalance:   balance.CurrentBalance,
			PendingBalance:   balance.PendingBalance,
			AvailableBalance: balance.AvailableBalance,
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

func (h *ApiHandlers) GetBalanceSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.BalanceUseCase.GetBalanceSummary(r.Context())
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	response := BalanceSummaryResponse{
		TotalAssets:      summary.TotalAssets,
		TotalLiabilities: summary.TotalLiabilities,
		NetWorth:         summary.NetWorth,
		LastCalculated:   summary.LastCalculated.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

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
