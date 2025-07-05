package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"finance/domain"
	"finance/domain/entities"
	"finance/internal/api/v1/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCreateExample(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockUC := &mocks.ExampleUseCaseMock{
			CreateExampleFunc: func(ctx context.Context, example entities.Example) (string, error) {
				return "123", nil
			},
		}

		h := &ApiHandlers{
			ExampleUseCase: mockUC,
		}

		reqBody := CreateExampleRequest{
			Title:   "Test Title",
			Content: "Test Content",
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/examples", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateExample(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response CreateExampleResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.ID != "123" {
			t.Errorf("expected ID '123', got '%s'", response.ID)
		}
	})

	t.Run("missing title", func(t *testing.T) {
		h := &ApiHandlers{
			ExampleUseCase: &mocks.ExampleUseCaseMock{},
		}

		reqBody := CreateExampleRequest{
			Content: "Test Content",
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/examples", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateExample(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("use case error", func(t *testing.T) {
		mockUC := &mocks.ExampleUseCaseMock{
			CreateExampleFunc: func(ctx context.Context, example entities.Example) (string, error) {
				return "", domain.ErrConflict
			},
		}

		h := &ApiHandlers{
			ExampleUseCase: mockUC,
		}

		reqBody := CreateExampleRequest{
			Title:   "Test Title",
			Content: "Test Content",
		}
		bodyJSON, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/examples", bytes.NewBuffer(bodyJSON))
		w := httptest.NewRecorder()

		h.CreateExample(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})
}

func TestGetExampleByID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		mockUC := &mocks.ExampleUseCaseMock{
			GetExampleByIDFunc: func(ctx context.Context, id string) (entities.Example, error) {
				return entities.Example{
					ID:      "123",
					Title:   "Test Title",
					Content: "Test Content",
				}, nil
			},
		}

		h := &ApiHandlers{
			ExampleUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/examples/123", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetExampleByID(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entities.Example
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.ID != "123" {
			t.Errorf("expected ID '123', got '%s'", response.ID)
		}
	})

	t.Run("missing id", func(t *testing.T) {
		h := &ApiHandlers{
			ExampleUseCase: &mocks.ExampleUseCaseMock{},
		}

		req := httptest.NewRequest(http.MethodGet, "/examples/", nil)
		w := httptest.NewRecorder()

		// Setup chi router context with empty id
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetExampleByID(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockUC := &mocks.ExampleUseCaseMock{
			GetExampleByIDFunc: func(ctx context.Context, id string) (entities.Example, error) {
				return entities.Example{}, domain.ErrNotFound
			},
		}

		h := &ApiHandlers{
			ExampleUseCase: mockUC,
		}

		req := httptest.NewRequest(http.MethodGet, "/examples/999", nil)
		w := httptest.NewRecorder()

		// Setup chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.GetExampleByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
