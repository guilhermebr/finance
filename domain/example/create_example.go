package example

import (
	"context"
	"fmt"
	"go-template/domain"
	"go-template/domain/entities"
)

func (uc UseCase) CreateExample(ctx context.Context, input entities.Example) (string, error) {
	if len(input.Title) == 0 {
		return "", fmt.Errorf("missing title: %w", domain.ErrMalformedParameters)
	}

	id, err := uc.R.CreateExample(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create example: %w", err)
	}

	return id, nil
}
