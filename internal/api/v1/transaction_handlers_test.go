package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"finance/domain/entities"
	"finance/internal/api/v1/mocks"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/guilhermebr/gox/monetary"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("successful creation with valid date", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
				// Return the transaction with generated fields
				monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(10050)) // $100.50
				return entities.Transaction{
					ID:          "test-123",
					AccountID:   transaction.AccountID,
					CategoryID:  transaction.CategoryID,
					Monetary:    *monetaryValue,
					Description: transaction.Description,
					Date:        transaction.Date,
					Status:      transaction.Status,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "100.50",
			Description: "Test transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}

		var response TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.ID != "test-123" {
			t.Errorf("expected ID 'test-123', got '%s'", response.ID)
		}
		if response.Date != "2024-01-15" {
			t.Errorf("expected date '2024-01-15', got '%s'", response.Date)
		}

		// Check if the usecase was called with the correct parsed date
		calls := mockUC.CreateTransactionCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 call to CreateTransaction, got %d", len(calls))
		}
		expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
		if !calls[0].Transaction.Date.Equal(expectedDate) {
			t.Errorf("expected date %v, got %v", expectedDate, calls[0].Transaction.Date)
		}
	})

	t.Run("successful creation with empty date (defaults to current date)", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
				monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(10050)) // $100.50
				return entities.Transaction{
					ID:          "test-123",
					AccountID:   transaction.AccountID,
					CategoryID:  transaction.CategoryID,
					Monetary:    *monetaryValue,
					Description: transaction.Description,
					Date:        transaction.Date,
					Status:      transaction.Status,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "100.50",
			Description: "Test transaction",
			Date:        "", // Empty date
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}

		// Check if the usecase was called with today's date
		calls := mockUC.CreateTransactionCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 call to CreateTransaction, got %d", len(calls))
		}

		// Should be roughly today (within a minute)
		now := time.Now()
		dateDiff := now.Sub(calls[0].Transaction.Date)
		if dateDiff > time.Minute || dateDiff < -time.Minute {
			t.Errorf("expected date to be close to now, got %v (diff: %v)", calls[0].Transaction.Date, dateDiff)
		}
	})

	t.Run("invalid date format", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "100.50",
			Description: "Test transaction",
			Date:        "invalid-date", // Invalid format
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response ErrorResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Error == "" {
			t.Error("expected error message, got empty string")
		}
		t.Logf("Error message: %s", response.Error)
	})

	t.Run("invalid date format - wrong separator", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "100.50",
			Description: "Test transaction",
			Date:        "2024/01/15", // Wrong separator
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response ErrorResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Error == "" {
			t.Error("expected error message, got empty string")
		}
		t.Logf("Error message: %s", response.Error)
	})

	t.Run("invalid date format - US format", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "100.50",
			Description: "Test transaction",
			Date:        "01/15/2024", // US format
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response ErrorResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Error == "" {
			t.Error("expected error message, got empty string")
		}
		t.Logf("Error message: %s", response.Error)
	})

	t.Run("invalid amount format", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "invalid-amount",
			Description: "Test transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response ErrorResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Error == "" {
			t.Error("expected error message, got empty string")
		}
		t.Logf("Error message: %s", response.Error)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte("{invalid json")))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			CreateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
				return entities.Transaction{}, errors.New("invalid data")
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		reqBody := CreateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "100.50",
			Description: "Test transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusPending,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetTransactionByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(10050)) // $100.50
		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionWithDetailsFunc: func(ctx context.Context, id string) (entities.Transaction, error) {
				return entities.Transaction{
					ID:          "test-123",
					AccountID:   "acc-1",
					CategoryID:  "cat-1",
					Monetary:    *monetaryValue,
					Description: "Test transaction",
					Date:        testDate,
					Status:      entities.TransactionStatusPending,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions/test-123", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetTransactionByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.ID != "test-123" {
			t.Errorf("expected ID 'test-123', got '%s'", response.ID)
		}
		if response.Date != "2024-01-15" {
			t.Errorf("expected date '2024-01-15', got '%s'", response.Date)
		}
	})

	t.Run("successful retrieval with related entities", func(t *testing.T) {
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(10050)) // $100.50
		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionWithDetailsFunc: func(ctx context.Context, id string) (entities.Transaction, error) {
				return entities.Transaction{
					ID:          "test-123",
					AccountID:   "acc-1",
					CategoryID:  "cat-1",
					Monetary:    *monetaryValue,
					Description: "Test transaction",
					Date:        testDate,
					Status:      entities.TransactionStatusPending,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Account: &entities.Account{
						ID:          "acc-1",
						Name:        "Test Account",
						Type:        entities.AccountTypeChecking,
						Description: "Test account",
					},
					Category: &entities.Category{
						ID:          "cat-1",
						Name:        "Test Category",
						Type:        entities.CategoryTypeExpense,
						Description: "Test category",
						Color:       "#FF0000",
					},
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions/test-123", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetTransactionByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.ID != "test-123" {
			t.Errorf("expected ID 'test-123', got '%s'", response.ID)
		}
		if response.Account == nil {
			t.Error("expected account to be populated")
		} else if response.Account.Name != "Test Account" {
			t.Errorf("expected account name 'Test Account', got '%s'", response.Account.Name)
		}
		if response.Category == nil {
			t.Error("expected category to be populated")
		} else if response.Category.Name != "Test Category" {
			t.Errorf("expected category name 'Test Category', got '%s'", response.Category.Name)
		}
	})

	t.Run("missing id", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions/", nil)
		w := httptest.NewRecorder()

		// Setup chi router context with empty id
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetTransactionByID(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionWithDetailsFunc: func(ctx context.Context, id string) (entities.Transaction, error) {
				return entities.Transaction{}, errors.New("database error")
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions/test-123", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetTransactionByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("transaction not found (empty ID)", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionWithDetailsFunc: func(ctx context.Context, id string) (entities.Transaction, error) {
				return entities.Transaction{}, nil // Empty transaction
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions/nonexistent", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetTransactionByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestGetAllTransactions(t *testing.T) {
	t.Run("successful retrieval with multiple transactions", func(t *testing.T) {
		testDate1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		testDate2 := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)
		monetaryValue1, err1 := monetary.NewMonetary(monetary.USD, big.NewInt(10050)) // $100.50
		if err1 != nil {
			t.Fatalf("Failed to create monetary value 1: %v", err1)
		}
		monetaryValue2, err2 := monetary.NewMonetary(monetary.USD, big.NewInt(5025)) // $50.25
		if err2 != nil {
			t.Fatalf("Failed to create monetary value 2: %v", err2)
		}

		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionsWithDetailsFunc: func(ctx context.Context, limit, offset int) ([]entities.Transaction, error) {
				return []entities.Transaction{
					{
						ID:          "test-123",
						AccountID:   "acc-1",
						CategoryID:  "cat-1",
						Monetary:    *monetaryValue1,
						Description: "Test transaction 1",
						Date:        testDate1,
						Status:      entities.TransactionStatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "test-456",
						AccountID:   "acc-2",
						CategoryID:  "cat-2",
						Monetary:    *monetaryValue2,
						Description: "Test transaction 2",
						Date:        testDate2,
						Status:      entities.TransactionStatusCleared,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		w := httptest.NewRecorder()

		h.GetAllTransactions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) != 2 {
			t.Errorf("expected 2 transactions, got %d", len(response))
		}

		// Check first transaction
		if response[0].ID != "test-123" {
			t.Errorf("expected first transaction ID 'test-123', got '%s'", response[0].ID)
		}
		if response[0].Date != "2024-01-15" {
			t.Errorf("expected first transaction date '2024-01-15', got '%s'", response[0].Date)
		}
		if response[0].Amount != "[USD ($) 100.50]" {
			t.Errorf("expected first transaction amount '[USD ($) 100.50]', got '%s'", response[0].Amount)
		}

		// Check second transaction
		if response[1].ID != "test-456" {
			t.Errorf("expected second transaction ID 'test-456', got '%s'", response[1].ID)
		}
		if response[1].Date != "2024-01-16" {
			t.Errorf("expected second transaction date '2024-01-16', got '%s'", response[1].Date)
		}
		if response[1].Amount != "[USD ($) 50.25]" {
			t.Errorf("expected second transaction amount '[USD ($) 50.25]', got '%s'", response[1].Amount)
		}

		// Check that usecase was called with correct parameters
		calls := mockUC.GetTransactionsWithDetailsCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 call to GetTransactionsWithDetails, got %d", len(calls))
		}
		if calls[0].Limit != 50 {
			t.Errorf("expected limit 50, got %d", calls[0].Limit)
		}
		if calls[0].Offset != 0 {
			t.Errorf("expected offset 0, got %d", calls[0].Offset)
		}
	})

	t.Run("successful retrieval with related entities", func(t *testing.T) {
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(10050)) // $100.50

		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionsWithDetailsFunc: func(ctx context.Context, limit, offset int) ([]entities.Transaction, error) {
				return []entities.Transaction{
					{
						ID:          "test-123",
						AccountID:   "acc-1",
						CategoryID:  "cat-1",
						Monetary:    *monetaryValue,
						Description: "Test transaction",
						Date:        testDate,
						Status:      entities.TransactionStatusPending,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Account: &entities.Account{
							ID:          "acc-1",
							Name:        "Test Account",
							Type:        entities.AccountTypeChecking,
							Description: "Test account",
						},
						Category: &entities.Category{
							ID:          "cat-1",
							Name:        "Test Category",
							Type:        entities.CategoryTypeExpense,
							Description: "Test category",
							Color:       "#FF0000",
						},
					},
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		w := httptest.NewRecorder()

		h.GetAllTransactions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) != 1 {
			t.Errorf("expected 1 transaction, got %d", len(response))
		}

		// Check that related entities are populated
		if response[0].Account == nil {
			t.Error("expected account to be populated")
		} else if response[0].Account.Name != "Test Account" {
			t.Errorf("expected account name 'Test Account', got '%s'", response[0].Account.Name)
		}
		if response[0].Category == nil {
			t.Error("expected category to be populated")
		} else if response[0].Category.Name != "Test Category" {
			t.Errorf("expected category name 'Test Category', got '%s'", response[0].Category.Name)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionsWithDetailsFunc: func(ctx context.Context, limit, offset int) ([]entities.Transaction, error) {
				return []entities.Transaction{}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		w := httptest.NewRecorder()

		h.GetAllTransactions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []TransactionResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) != 0 {
			t.Errorf("expected 0 transactions, got %d", len(response))
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			GetTransactionsWithDetailsFunc: func(ctx context.Context, limit, offset int) ([]entities.Transaction, error) {
				return nil, errors.New("database error")
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		w := httptest.NewRecorder()

		h.GetAllTransactions(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestUpdateTransaction(t *testing.T) {
	t.Run("successful update with valid date", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			UpdateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
				monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(15075)) // $150.75
				return entities.Transaction{
					ID:          transaction.ID,
					AccountID:   transaction.AccountID,
					CategoryID:  transaction.CategoryID,
					Monetary:    *monetaryValue,
					Description: transaction.Description,
					Date:        transaction.Date,
					Status:      transaction.Status,
					CreatedAt:   time.Now().Add(-time.Hour),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		reqBody := UpdateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "150.75",
			Description: "Updated transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusCleared,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/transactions/test-123", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}

		// Check if the usecase was called with the correct parsed date
		calls := mockUC.UpdateTransactionCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 call to UpdateTransaction, got %d", len(calls))
		}
		expectedDate, _ := time.Parse("2006-01-02", "2024-01-15")
		if !calls[0].Transaction.Date.Equal(expectedDate) {
			t.Errorf("expected date %v, got %v", expectedDate, calls[0].Transaction.Date)
		}
		if calls[0].Transaction.ID != "test-123" {
			t.Errorf("expected transaction ID 'test-123', got '%s'", calls[0].Transaction.ID)
		}
	})

	t.Run("successful update with empty date", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			UpdateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
				monetaryValue, _ := monetary.NewMonetary(monetary.USD, big.NewInt(15075)) // $150.75
				return entities.Transaction{
					ID:          transaction.ID,
					AccountID:   transaction.AccountID,
					CategoryID:  transaction.CategoryID,
					Monetary:    *monetaryValue,
					Description: transaction.Description,
					Date:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), // Existing date
					Status:      transaction.Status,
					CreatedAt:   time.Now().Add(-time.Hour),
					UpdatedAt:   time.Now(),
				}, nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		reqBody := UpdateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "150.75",
			Description: "Updated transaction",
			Date:        "", // Empty date - should preserve existing
			Status:      entities.TransactionStatusCleared,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/transactions/test-123", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}

		// Check if the usecase was called with zero time (letting usecase handle it)
		calls := mockUC.UpdateTransactionCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 call to UpdateTransaction, got %d", len(calls))
		}
		if !calls[0].Transaction.Date.IsZero() {
			t.Errorf("expected zero date for empty date input, got %v", calls[0].Transaction.Date)
		}
	})

	t.Run("invalid date format in update", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := UpdateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "150.75",
			Description: "Updated transaction",
			Date:        "invalid-date",
			Status:      entities.TransactionStatusCleared,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/transactions/test-123", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response ErrorResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Error == "" {
			t.Error("expected error message, got empty string")
		}
		t.Logf("Error message: %s", response.Error)
	})

	t.Run("invalid amount format in update", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := UpdateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "invalid-amount",
			Description: "Updated transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusCleared,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/transactions/test-123", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response ErrorResponseBody
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Error == "" {
			t.Error("expected error message, got empty string")
		}
		t.Logf("Error message: %s", response.Error)
	})

	t.Run("missing id", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		reqBody := UpdateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "150.75",
			Description: "Updated transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusCleared,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/transactions/", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		// Setup chi router context with empty id
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		req := httptest.NewRequest(http.MethodPut, "/transactions/test-123", bytes.NewBuffer([]byte("{invalid json")))
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			UpdateTransactionFunc: func(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
				return entities.Transaction{}, errors.New("validation error")
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		reqBody := UpdateTransactionRequest{
			AccountID:   "acc-1",
			CategoryID:  "cat-1",
			Amount:      "150.75",
			Description: "Updated transaction",
			Date:        "2024-01-15",
			Status:      entities.TransactionStatusCleared,
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/transactions/test-123", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.UpdateTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestDeleteTransaction(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			DeleteTransactionFunc: func(ctx context.Context, id string) error {
				return nil
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodDelete, "/transactions/test-123", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.DeleteTransaction(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		// Check that usecase was called with correct ID
		calls := mockUC.DeleteTransactionCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 call to DeleteTransaction, got %d", len(calls))
		}
		if calls[0].ID != "test-123" {
			t.Errorf("expected transaction ID 'test-123', got '%s'", calls[0].ID)
		}
	})

	t.Run("missing id", func(t *testing.T) {
		h := &ApiHandlers{
			TransactionUseCase: &mocks.TransactionUseCaseMock{},
		}

		req := httptest.NewRequest(http.MethodDelete, "/transactions/", nil)
		w := httptest.NewRecorder()

		// Setup chi router context with empty id
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.DeleteTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUC := &mocks.TransactionUseCaseMock{
			DeleteTransactionFunc: func(ctx context.Context, id string) error {
				return errors.New("transaction not found")
			},
		}

		h := &ApiHandlers{
			TransactionUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodDelete, "/transactions/nonexistent", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.DeleteTransaction(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
