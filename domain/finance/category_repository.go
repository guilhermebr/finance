package finance

import (
	"context"
	"finance/domain/entities"
)

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/category_repository.go . CategoryRepository
type CategoryRepository interface {
	CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error)
	GetCategoryByID(ctx context.Context, id string) (entities.Category, error)
	GetAllCategories(ctx context.Context) ([]entities.Category, error)
	GetCategoriesByType(ctx context.Context, categoryType entities.CategoryType) ([]entities.Category, error)
	UpdateCategory(ctx context.Context, category entities.Category) (entities.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}
