package pg

import (
	"context"
	"database/sql"
	"finance/domain/entities"
	"finance/internal/repository/pg/gen"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	queries *gen.Queries
	db      *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		queries: gen.New(db),
		db:      db,
	}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {
	result, err := r.queries.CreateCategory(ctx, category.Name, string(category.Type), category.Description, category.Color)
	if err != nil {
		return entities.Category{}, err
	}

	return entities.Category{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.CategoryType(result.Type),
		Description: result.Description,
		Color:       result.Color,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id string) (entities.Category, error) {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return entities.Category{}, err
	}

	result, err := r.queries.GetCategoryByID(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Category{}, nil
		}
		return entities.Category{}, err
	}

	return entities.Category{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.CategoryType(result.Type),
		Description: result.Description,
		Color:       result.Color,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	results, err := r.queries.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]entities.Category, len(results))
	for i, result := range results {
		categories[i] = entities.Category{
			ID:          result.ID.String(),
			Name:        result.Name,
			Type:        entities.CategoryType(result.Type),
			Description: result.Description,
			Color:       result.Color,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoriesByType(ctx context.Context, categoryType entities.CategoryType) ([]entities.Category, error) {
	results, err := r.queries.GetCategoriesByType(ctx, string(categoryType))
	if err != nil {
		return nil, err
	}

	categories := make([]entities.Category, len(results))
	for i, result := range results {
		categories[i] = entities.Category{
			ID:          result.ID.String(),
			Name:        result.Name,
			Type:        entities.CategoryType(result.Type),
			Description: result.Description,
			Color:       result.Color,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {
	uuid, err := uuid.FromString(category.ID)
	if err != nil {
		return entities.Category{}, err
	}

	result, err := r.queries.UpdateCategory(ctx, uuid, category.Name, string(category.Type), category.Description, category.Color)
	if err != nil {
		return entities.Category{}, err
	}

	return entities.Category{
		ID:          result.ID.String(),
		Name:        result.Name,
		Type:        entities.CategoryType(result.Type),
		Description: result.Description,
		Color:       result.Color,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	uuid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteCategory(ctx, uuid)
}
