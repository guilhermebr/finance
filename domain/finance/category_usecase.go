package finance

import (
	"context"
	"finance/domain/entities"
	"fmt"
	"strings"
)

type CategoryUseCase struct {
	categoryRepo CategoryRepository
}

func NewCategoryUseCase(categoryRepo CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {
	// Validate input
	if err := uc.validateCategory(category); err != nil {
		return entities.Category{}, err
	}

	// Set default color if not provided
	if category.Color == "" {
		category.Color = "#6B7280" // Default gray color
	}

	createdCategory, err := uc.categoryRepo.CreateCategory(ctx, category)
	if err != nil {
		return entities.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return createdCategory, nil
}

func (uc *CategoryUseCase) GetCategoryByID(ctx context.Context, id string) (entities.Category, error) {
	if id == "" {
		return entities.Category{}, fmt.Errorf("category ID cannot be empty")
	}

	category, err := uc.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return entities.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

func (uc *CategoryUseCase) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	categories, err := uc.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (uc *CategoryUseCase) GetCategoriesByType(ctx context.Context, categoryType entities.CategoryType) ([]entities.Category, error) {
	if categoryType == "" {
		return nil, fmt.Errorf("category type cannot be empty")
	}

	categories, err := uc.categoryRepo.GetCategoriesByType(ctx, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by type: %w", err)
	}

	return categories, nil
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {
	// Validate input
	if err := uc.validateCategory(category); err != nil {
		return entities.Category{}, err
	}

	if category.ID == "" {
		return entities.Category{}, fmt.Errorf("category ID cannot be empty")
	}

	// Check if category exists
	existingCategory, err := uc.categoryRepo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return entities.Category{}, fmt.Errorf("failed to get existing category: %w", err)
	}

	if existingCategory.ID == "" {
		return entities.Category{}, fmt.Errorf("category not found")
	}

	// Set default color if not provided
	if category.Color == "" {
		category.Color = "#6B7280" // Default gray color
	}

	updatedCategory, err := uc.categoryRepo.UpdateCategory(ctx, category)
	if err != nil {
		return entities.Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	return updatedCategory, nil
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("category ID cannot be empty")
	}

	// Check if category exists
	category, err := uc.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	if category.ID == "" {
		return fmt.Errorf("category not found")
	}

	// TODO: Check if category is being used by transactions
	// For now, we'll rely on the database constraint to prevent deletion

	err = uc.categoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (uc *CategoryUseCase) validateCategory(category entities.Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return fmt.Errorf("category name cannot be empty")
	}

	if category.Type == "" {
		return fmt.Errorf("category type cannot be empty")
	}

	validTypes := []entities.CategoryType{
		entities.CategoryTypeIncome,
		entities.CategoryTypeExpense,
	}

	isValidType := false
	for _, validType := range validTypes {
		if category.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid category type: %s", category.Type)
	}

	return nil
}
