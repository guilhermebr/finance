package v1

import (
	"context"
	"encoding/json"
	"errors"
	"go-template/domain"
	"go-template/domain/entities"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CreateExampleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateExampleResponse struct {
	ID string `json:"id"`
}

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/example_uc.go . ExampleUseCase
type ExampleUseCase interface {
	CreateExample(ctx context.Context, example entities.Example) (string, error)
	GetExampleByID(ctx context.Context, id string) (entities.Example, error)
}

// CreateExample godoc
// @Summary Create a new example
// @Description Create a new example with the given title and content
// @Tags examples
// @Accept json
// @Produce json
// @Param example body CreateExampleRequest true "Example to create"
func (h *ApiHandlers) CreateExample(w http.ResponseWriter, r *http.Request) {
	var input CreateExampleRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		errorResponse(w, r, http.StatusBadRequest, errors.New("title is required"))
		return
	}

	example := entities.Example{
		Title:   input.Title,
		Content: input.Content,
	}

	id, err := h.ExampleUseCase.CreateExample(r.Context(), example)
	if err != nil {
		slog.Error("failed to create example", "error", err, "input", input)
		switch {
		case errors.Is(err, domain.ErrMalformedParameters):
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		case errors.Is(err, domain.ErrConflict):
			errorResponse(w, r, http.StatusConflict, err)
			return
		case errors.Is(err, domain.ErrDuplicateKey):
			errorResponse(w, r, http.StatusConflict, err)
			return
		default:
			unknownErrorResponse(w, r)
			return
		}
	}

	slog.Info("example created successfully", "id", id)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateExampleResponse{ID: id})
}

// GetExampleByID godoc
// @Summary Get an example by ID
// @Description Get an example by its unique identifier
// @Tags examples
// @Accept json
// @Produce json
// @Param id path string true "Example ID"
func (h *ApiHandlers) GetExampleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	example, err := h.ExampleUseCase.GetExampleByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get example", "error", err, "id", id)
		switch {
		case errors.Is(err, domain.ErrNotFound):
			errorResponse(w, r, http.StatusNotFound, errors.New("example not found"))
			return
		case errors.Is(err, domain.ErrMalformedParameters):
			errorResponse(w, r, http.StatusBadRequest, err)
			return
		default:
			unknownErrorResponse(w, r)
			return
		}
	}

	if example.ID == "" {
		errorResponse(w, r, http.StatusNotFound, errors.New("example not found"))
		return
	}

	slog.Info("example retrieved successfully", "id", id)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, example)
}
