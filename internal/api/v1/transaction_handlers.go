package v1

import (
	"context"
	"encoding/json"
	"finance/domain/entities"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"math/big"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/guilhermebr/gox/monetary"
)

// Transaction request/response types
type CreateTransactionRequest struct {
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      string                     `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
}

type UpdateTransactionRequest struct {
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      string                     `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
}

type TransactionResponse struct {
	ID          string                     `json:"id"`
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      string                     `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
	CreatedAt   string                     `json:"created_at"`
	UpdatedAt   string                     `json:"updated_at"`
	Account     *AccountResponse           `json:"account,omitempty"`
	Category    *CategoryResponse          `json:"category,omitempty"`
}

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/transaction_uc.go . TransactionUseCase
type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error)
	GetTransactionWithDetails(ctx context.Context, id string) (entities.Transaction, error)
	GetTransactionsWithDetails(ctx context.Context, limit int, offset int) ([]entities.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error)
	DeleteTransaction(ctx context.Context, id string) error
}

// Transaction handlers

// CreateTransaction creates a new transaction
//
//	@Summary		Create a new transaction
//	@Description	Create a new financial transaction with the provided details
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			transaction	body		CreateTransactionRequest	true	"Transaction data"
//	@Success		201			{object}	TransactionResponse			"Transaction created successfully"
//	@Failure		400			{object}	ErrorResponseBody			"Bad request"
//	@Router			/transactions [post]
func (h *ApiHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode transaction request", "error", err)
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Parse date - default to current date if empty
	var transactionDate time.Time
	if req.Date != "" {
		var err error
		transactionDate, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			slog.Error("failed to parse date request", "error", err, "date", req.Date)
			errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("date", "must be in format YYYY-MM-DD"))
			return
		}
	} else {
		// Default to current date if no date provided
		transactionDate = time.Now()
	}

	// Parse amount as float and create temporary monetary value with USD
	// The use case will handle the proper asset conversion based on the account
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		slog.Error("failed to parse amount", "error", err, "amount", req.Amount)
		errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("amount", "must be a valid decimal number"))
		return
	}

	// Create temporary monetary value with USD - use case will convert to correct asset
	tempMonetary, err := monetary.NewMonetary(monetary.USD, big.NewInt(int64(amountFloat*100)))
	if err != nil {
		slog.Error("failed to create monetary value", "error", err, "amount", req.Amount)
		errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("amount", "must be a valid decimal number"))
		return
	}

	// Create transaction entity
	transaction := entities.Transaction{
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Monetary:    *tempMonetary,
		Description: req.Description,
		Date:        transactionDate,
		Status:      req.Status,
	}

	createdTransaction, err := h.TransactionUseCase.CreateTransaction(r.Context(), transaction)
	if err != nil {
		slog.Error("failed to create transaction", "error", err, "account_id", req.AccountID, "category_id", req.CategoryID, "amount", req.Amount)
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := TransactionResponse{
		ID:          createdTransaction.ID,
		AccountID:   createdTransaction.AccountID,
		CategoryID:  createdTransaction.CategoryID,
		Amount:      createdTransaction.Monetary.String(),
		Description: createdTransaction.Description,
		Date:        createdTransaction.Date.Format("2006-01-02"),
		Status:      createdTransaction.Status,
		CreatedAt:   createdTransaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   createdTransaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// GetTransactionByID retrieves a transaction by its ID
//
//	@Summary		Get transaction by ID
//	@Description	Retrieve a specific transaction by its unique identifier
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Transaction ID"
//	@Success		200	{object}	TransactionResponse	"Transaction retrieved successfully"
//	@Failure		400	{object}	ErrorResponseBody	"Bad request"
//	@Failure		404	{object}	ErrorResponseBody	"Transaction not found"
//	@Router			/transactions/{id} [get]
func (h *ApiHandlers) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("missing transaction ID parameter")
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	transaction, err := h.TransactionUseCase.GetTransactionWithDetails(r.Context(), id)
	if err != nil {
		slog.Error("failed to get transaction", "error", err, "transaction_id", id)
		errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	if transaction.ID == "" {
		slog.Error("transaction not found", "transaction_id", id)
		errorResponse(w, r, http.StatusNotFound, errNotFound("transaction"))
		return
	}

	response := TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Monetary.String(),
		Description: transaction.Description,
		Date:        transaction.Date.Format("2006-01-02"),
		Status:      transaction.Status,
		CreatedAt:   transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add related entities if available
	if transaction.Account != nil {
		response.Account = &AccountResponse{
			ID:          transaction.Account.ID,
			Name:        transaction.Account.Name,
			Type:        transaction.Account.Type,
			Asset:       transaction.Account.Asset.Asset,
			Description: transaction.Account.Description,
		}
	}

	if transaction.Category != nil {
		response.Category = &CategoryResponse{
			ID:          transaction.Category.ID,
			Name:        transaction.Category.Name,
			Type:        transaction.Category.Type,
			Description: transaction.Category.Description,
			Color:       transaction.Category.Color,
		}
	}

	render.JSON(w, r, response)
}

// GetAllTransactions retrieves all transactions
//
//	@Summary		Get all transactions
//	@Description	Retrieve a list of all financial transactions with pagination (limit: 50, offset: 0)
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		TransactionResponse	"Transactions retrieved successfully"
//	@Failure		500	{object}	ErrorResponseBody	"Internal server error"
//	@Router			/transactions [get]
func (h *ApiHandlers) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.TransactionUseCase.GetTransactionsWithDetails(r.Context(), 50, 0)
	if err != nil {
		slog.Error("failed to get transactions", "error", err)
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = TransactionResponse{
			ID:          transaction.ID,
			AccountID:   transaction.AccountID,
			CategoryID:  transaction.CategoryID,
			Amount:      transaction.Monetary.String(),
			Description: transaction.Description,
			Date:        transaction.Date.Format("2006-01-02"),
			Status:      transaction.Status,
			CreatedAt:   transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Add related entities if available
		if transaction.Account != nil {
			responses[i].Account = &AccountResponse{
				ID:          transaction.Account.ID,
				Name:        transaction.Account.Name,
				Type:        transaction.Account.Type,
				Asset:       transaction.Account.Asset.Asset,
				Description: transaction.Account.Description,
			}
		}

		if transaction.Category != nil {
			responses[i].Category = &CategoryResponse{
				ID:          transaction.Category.ID,
				Name:        transaction.Category.Name,
				Type:        transaction.Category.Type,
				Description: transaction.Category.Description,
				Color:       transaction.Category.Color,
			}
		}
	}

	render.JSON(w, r, responses)
}

// UpdateTransaction updates an existing transaction
//
//	@Summary		Update transaction
//	@Description	Update an existing transaction with new information
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string						true	"Transaction ID"
//	@Param			transaction	body		UpdateTransactionRequest	true	"Updated transaction data"
//	@Success		200			{object}	TransactionResponse			"Transaction updated successfully"
//	@Failure		400			{object}	ErrorResponseBody			"Bad request"
//	@Failure		404			{object}	ErrorResponseBody			"Transaction not found"
//	@Router			/transactions/{id} [put]
func (h *ApiHandlers) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("missing transaction ID parameter")
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode update transaction request", "error", err, "transaction_id", id)
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Parse amount as float and create temporary monetary value with USD
	// The use case will handle the proper asset conversion based on the account
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		slog.Error("failed to parse amount", "error", err, "amount", req.Amount, "transaction_id", id)
		errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("amount", "must be a valid decimal number"))
		return
	}

	// Create temporary monetary value with USD - use case will convert to correct asset
	tempMonetary, err := monetary.NewMonetary(monetary.USD, big.NewInt(int64(amountFloat*100)))
	if err != nil {
		slog.Error("failed to create monetary value", "error", err, "amount", req.Amount, "transaction_id", id)
		errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("amount", "must be a valid decimal number"))
		return
	}

	var transactionDate time.Time
	if req.Date != "" {
		var err error
		transactionDate, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			slog.Error("failed to parse date request", "error", err, "date", req.Date, "transaction_id", id)
			errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("date", "must be in format YYYY-MM-DD"))
			return
		}
	}

	transaction := entities.Transaction{
		ID:          id,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Monetary:    *tempMonetary,
		Description: req.Description,
		Date:        transactionDate,
		Status:      req.Status,
	}

	updatedTransaction, err := h.TransactionUseCase.UpdateTransaction(r.Context(), transaction)
	if err != nil {
		slog.Error("failed to update transaction", "error", err, "transaction_id", id, "account_id", req.AccountID, "category_id", req.CategoryID)
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := TransactionResponse{
		ID:          updatedTransaction.ID,
		AccountID:   updatedTransaction.AccountID,
		CategoryID:  updatedTransaction.CategoryID,
		Amount:      updatedTransaction.Monetary.String(),
		Description: updatedTransaction.Description,
		Date:        updatedTransaction.Date.Format("2006-01-02"),
		Status:      updatedTransaction.Status,
		CreatedAt:   updatedTransaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updatedTransaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

// DeleteTransaction deletes a transaction
//
//	@Summary		Delete transaction
//	@Description	Delete a transaction by its ID
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Transaction ID"
//	@Success		204	"Transaction deleted successfully"
//	@Failure		400	{object}	ErrorResponseBody	"Bad request"
//	@Failure		404	{object}	ErrorResponseBody	"Transaction not found"
//	@Router			/transactions/{id} [delete]
func (h *ApiHandlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		slog.Error("missing transaction ID parameter")
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	err := h.TransactionUseCase.DeleteTransaction(r.Context(), id)
	if err != nil {
		slog.Error("failed to delete transaction", "error", err, "transaction_id", id)
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
