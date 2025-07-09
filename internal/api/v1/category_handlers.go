package v1

import (
	"context"
	"encoding/json"
	"finance/domain/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Category request/response types
type CreateCategoryRequest struct {
	Name        string                `json:"name"`
	Type        entities.CategoryType `json:"type"`
	Description string                `json:"description"`
	Color       string                `json:"color"`
}

type UpdateCategoryRequest struct {
	Name        string                `json:"name"`
	Type        entities.CategoryType `json:"type"`
	Description string                `json:"description"`
	Color       string                `json:"color"`
}

type CategoryResponse struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Type        entities.CategoryType `json:"type"`
	Description string                `json:"description"`
	Color       string                `json:"color"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
}

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/category_uc.go . CategoryUseCase
type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error)
	GetCategoryByID(ctx context.Context, id string) (entities.Category, error)
	GetAllCategories(ctx context.Context) ([]entities.Category, error)
	UpdateCategory(ctx context.Context, category entities.Category) (entities.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

// Category handlers

// CreateCategory creates a new category
//
//	@Summary		Create a new category
//	@Description	Create a new transaction category with the provided details
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			category	body		CreateCategoryRequest	true	"Category data"
//	@Success		201			{object}	CategoryResponse		"Category created successfully"
//	@Failure		400			{object}	ErrorResponseBody		"Bad request"
//	@Router			/categories [post]
func (h *ApiHandlers) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	category := entities.Category{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Color:       req.Color,
	}

	createdCategory, err := h.CategoryUseCase.CreateCategory(r.Context(), category)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := CategoryResponse{
		ID:          createdCategory.ID,
		Name:        createdCategory.Name,
		Type:        createdCategory.Type,
		Description: createdCategory.Description,
		Color:       createdCategory.Color,
		CreatedAt:   createdCategory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   createdCategory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

// GetCategoryByID retrieves a category by its ID
//
//	@Summary		Get category by ID
//	@Description	Retrieve a specific category by its unique identifier
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Category ID"
//	@Success		200	{object}	CategoryResponse	"Category retrieved successfully"
//	@Failure		400	{object}	ErrorResponseBody	"Bad request"
//	@Failure		404	{object}	ErrorResponseBody	"Category not found"
//	@Router			/categories/{id} [get]
func (h *ApiHandlers) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	category, err := h.CategoryUseCase.GetCategoryByID(r.Context(), id)
	if err != nil {
		errorResponse(w, r, http.StatusNotFound, err)
		return
	}

	if category.ID == "" {
		errorResponse(w, r, http.StatusNotFound, errNotFound("category"))
		return
	}

	response := CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Type:        category.Type,
		Description: category.Description,
		Color:       category.Color,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

// GetAllCategories retrieves all categories
//
//	@Summary		Get all categories
//	@Description	Retrieve a list of all transaction categories
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		CategoryResponse	"Categories retrieved successfully"
//	@Failure		500	{object}	ErrorResponseBody	"Internal server error"
//	@Router			/categories [get]
func (h *ApiHandlers) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.CategoryUseCase.GetAllCategories(r.Context())
	if err != nil {
		errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Type:        category.Type,
			Description: category.Description,
			Color:       category.Color,
			CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	render.JSON(w, r, responses)
}

// UpdateCategory updates an existing category
//
//	@Summary		Update category
//	@Description	Update an existing category with new information
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"Category ID"
//	@Param			category	body		UpdateCategoryRequest	true	"Updated category data"
//	@Success		200			{object}	CategoryResponse		"Category updated successfully"
//	@Failure		400			{object}	ErrorResponseBody		"Bad request"
//	@Failure		404			{object}	ErrorResponseBody		"Category not found"
//	@Router			/categories/{id} [put]
func (h *ApiHandlers) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	category := entities.Category{
		ID:          id,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Color:       req.Color,
	}

	updatedCategory, err := h.CategoryUseCase.UpdateCategory(r.Context(), category)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	response := CategoryResponse{
		ID:          updatedCategory.ID,
		Name:        updatedCategory.Name,
		Type:        updatedCategory.Type,
		Description: updatedCategory.Description,
		Color:       updatedCategory.Color,
		CreatedAt:   updatedCategory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updatedCategory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

// DeleteCategory deletes a category
//
//	@Summary		Delete category
//	@Description	Delete a category by its ID
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Category ID"
//	@Success		204	"Category deleted successfully"
//	@Failure		400	{object}	ErrorResponseBody	"Bad request"
//	@Failure		404	{object}	ErrorResponseBody	"Category not found"
//	@Router			/categories/{id} [delete]
func (h *ApiHandlers) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		errorResponse(w, r, http.StatusBadRequest, errMissingParameter("id"))
		return
	}

	err := h.CategoryUseCase.DeleteCategory(r.Context(), id)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
