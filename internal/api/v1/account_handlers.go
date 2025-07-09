package v1

import (
	"context"
	"encoding/json"
	"finance/domain/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/guilhermebr/gox/monetary"
)

// Account request/response types
type CreateAccountRequest struct {
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Asset       string               `json:"asset"`
	Description string               `json:"description"`
}

type UpdateAccountRequest struct {
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Asset       string               `json:"asset"`
	Description string               `json:"description"`
}

type AccountResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Type        entities.AccountType `json:"type"`
	Asset       string               `json:"asset"`
	Description string               `json:"description"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/account_uc.go . AccountUseCase
type AccountUseCase interface {
	CreateAccount(ctx context.Context, account entities.Account) (entities.Account, error)
	GetAccountByID(ctx context.Context, id string) (entities.Account, error)
	GetAllAccounts(ctx context.Context) ([]entities.Account, error)
	UpdateAccount(ctx context.Context, account entities.Account) (entities.Account, error)
	DeleteAccount(ctx context.Context, id string) error
}

// Account handlers

// CreateAccount creates a new account
//
//	@Summary		Create a new account
//	@Description	Create a new financial account with the provided details
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			account	body		CreateAccountRequest	true	"Account data"
//	@Success		201		{object}	AccountResponse			"Account created successfully"
//	@Failure		400		{object}	ErrorResponseBody		"Bad request"
//	@Router			/accounts [post]
func (h *ApiHandlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Validate and get asset
	asset, ok := monetary.FindAssetByName(req.Asset)
	if !ok {
		errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("asset", req.Asset))
		return
	}

	account := entities.Account{
		Name:        req.Name,
		Type:        req.Type,
		Asset:       asset,
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
		Asset:       createdAccount.Asset.Asset,
		Description: createdAccount.Description,
		CreatedAt:   createdAccount.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   createdAccount.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// GetAccountByID retrieves an account by its ID
//
//	@Summary		Get account by ID
//	@Description	Retrieve a specific account by its unique identifier
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Account ID"
//	@Success		200	{object}	AccountResponse	"Account retrieved successfully"
//	@Failure		400	{object}	ErrorResponseBody	"Bad request"
//	@Failure		404	{object}	ErrorResponseBody	"Account not found"
//	@Router			/accounts/{id} [get]
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
		Asset:       account.Asset.Asset,
		Description: account.Description,
		CreatedAt:   account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

// GetAllAccounts retrieves all accounts
//
//	@Summary		Get all accounts
//	@Description	Retrieve a list of all financial accounts
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		AccountResponse		"Accounts retrieved successfully"
//	@Failure		500	{object}	ErrorResponseBody	"Internal server error"
//	@Router			/accounts [get]
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
			Asset:       account.Asset.Asset,
			Description: account.Description,
			CreatedAt:   account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	render.JSON(w, r, responses)
}

// UpdateAccount updates an existing account
//
//	@Summary		Update account
//	@Description	Update an existing account with new information
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Account ID"
//	@Param			account	body		UpdateAccountRequest	true	"Updated account data"
//	@Success		200		{object}	AccountResponse		"Account updated successfully"
//	@Failure		400		{object}	ErrorResponseBody	"Bad request"
//	@Failure		404		{object}	ErrorResponseBody	"Account not found"
//	@Router			/accounts/{id} [put]
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

	// Validate and get asset
	asset, ok := monetary.FindAssetByName(req.Asset)
	if !ok {
		errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("asset", req.Asset))
		return
	}

	account := entities.Account{
		ID:          id,
		Name:        req.Name,
		Type:        req.Type,
		Asset:       asset,
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
		Asset:       updatedAccount.Asset.Asset,
		Description: updatedAccount.Description,
		CreatedAt:   updatedAccount.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updatedAccount.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

// DeleteAccount deletes an account
//
//	@Summary		Delete account
//	@Description	Delete an account by its ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Account ID"
//	@Success		204	"Account deleted successfully"
//	@Failure		400	{object}	ErrorResponseBody	"Bad request"
//	@Failure		404	{object}	ErrorResponseBody	"Account not found"
//	@Router			/accounts/{id} [delete]
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
