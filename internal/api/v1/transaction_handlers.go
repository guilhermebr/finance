package v1

import (
	"encoding/json"
	"finance/domain/entities"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Transaction request/response types
type CreateTransactionRequest struct {
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      float64                    `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
}

type UpdateTransactionRequest struct {
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      float64                    `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
}

type TransactionResponse struct {
	ID          string                     `json:"id"`
	AccountID   string                     `json:"account_id"`
	CategoryID  string                     `json:"category_id"`
	Amount      float64                    `json:"amount"`
	Description string                     `json:"description"`
	Date        string                     `json:"date"`
	Status      entities.TransactionStatus `json:"status"`
	CreatedAt   string                     `json:"created_at"`
	UpdatedAt   string                     `json:"updated_at"`
	Account     *AccountResponse           `json:"account,omitempty"`
	Category    *CategoryResponse          `json:"category,omitempty"`
}

// Transaction handlers
func (h *ApiHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Parse date
	var transactionDate time.Time
	if req.Date != "" {
		var err error
		transactionDate, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("date", req.Date))
			return
		}
	}

	transaction := entities.Transaction{
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        transactionDate,
		Status:      req.Status,
	}

	createdTransaction, err := h.TransactionUseCase.CreateTransaction(r.Context(), transaction)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := TransactionResponse{
		ID:          createdTransaction.ID,
		AccountID:   createdTransaction.AccountID,
		CategoryID:  createdTransaction.CategoryID,
		Amount:      createdTransaction.Amount,
		Description: createdTransaction.Description,
		Date:        createdTransaction.Date.Format("2006-01-02"),
		Status:      createdTransaction.Status,
		CreatedAt:   createdTransaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   createdTransaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

func (h *ApiHandlers) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	transaction, err := h.TransactionUseCase.GetTransactionWithDetails(r.Context(), id)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	if transaction.ID == "" {
		errorResponse(w, r, http.StatusNotFound, errNotFound("transaction"))
		return
	}

	response := TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		CategoryID:  transaction.CategoryID,
		Amount:      transaction.Amount,
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

func (h *ApiHandlers) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.TransactionUseCase.GetTransactionsWithDetails(r.Context(), 50, 0)
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = TransactionResponse{
			ID:          transaction.ID,
			AccountID:   transaction.AccountID,
			CategoryID:  transaction.CategoryID,
			Amount:      transaction.Amount,
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

func (h *ApiHandlers) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// Parse date
	var transactionDate time.Time
	if req.Date != "" {
		var err error
		transactionDate, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			errorResponse(w, r, http.StatusBadRequest, errInvalidParameter("date", req.Date))
			return
		}
	}

	transaction := entities.Transaction{
		ID:          id,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        transactionDate,
		Status:      req.Status,
	}

	updatedTransaction, err := h.TransactionUseCase.UpdateTransaction(r.Context(), transaction)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := TransactionResponse{
		ID:          updatedTransaction.ID,
		AccountID:   updatedTransaction.AccountID,
		CategoryID:  updatedTransaction.CategoryID,
		Amount:      updatedTransaction.Amount,
		Description: updatedTransaction.Description,
		Date:        updatedTransaction.Date.Format("2006-01-02"),
		Status:      updatedTransaction.Status,
		CreatedAt:   updatedTransaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updatedTransaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

func (h *ApiHandlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	err := h.TransactionUseCase.DeleteTransaction(r.Context(), id)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
