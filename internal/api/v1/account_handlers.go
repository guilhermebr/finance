package v1

import (
	"encoding/json"
	"finance/domain/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Account request/response types
type CreateAccountRequest struct {
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Description string               `json:"description"`
}

type UpdateAccountRequest struct {
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Description string               `json:"description"`
}

type AccountResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Description string               `json:"description"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

// Account handlers
func (h *ApiHandlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	account := entities.Account{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	}

	createdAccount, err := h.AccountUseCase.CreateAccount(r.Context(), account)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := AccountResponse{
		ID:          createdAccount.ID,
		Name:        createdAccount.Name,
		Type:        createdAccount.Type,
		Description: createdAccount.Description,
		CreatedAt:   createdAccount.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   createdAccount.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

func (h *ApiHandlers) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	account, err := h.AccountUseCase.GetAccountByID(r.Context(), id)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	if account.ID == "" {
		errorResponse(w, r, http.StatusNotFound, errNotFound("account"))
		return
	}

	response := AccountResponse{
		ID:          account.ID,
		Name:        account.Name,
		Type:        account.Type,
		Description: account.Description,
		CreatedAt:   account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

func (h *ApiHandlers) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.AccountUseCase.GetAllAccounts(r.Context())
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	responses := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = AccountResponse{
			ID:          account.ID,
			Name:        account.Name,
			Type:        account.Type,
			Description: account.Description,
			CreatedAt:   account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	render.JSON(w, r, responses)
}

func (h *ApiHandlers) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	var req UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	account := entities.Account{
		ID:          id,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	}

	updatedAccount, err := h.AccountUseCase.UpdateAccount(r.Context(), account)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := AccountResponse{
		ID:          updatedAccount.ID,
		Name:        updatedAccount.Name,
		Type:        updatedAccount.Type,
		Description: updatedAccount.Description,
		CreatedAt:   updatedAccount.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updatedAccount.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

func (h *ApiHandlers) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	err := h.AccountUseCase.DeleteAccount(r.Context(), id)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
