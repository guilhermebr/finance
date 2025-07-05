package example

import (
	"context"
	"fmt"
	"finance/domain"
	"finance/domain/entities"
)

func (uc UseCase) GetExampleByID(ctx context.Context, id string) (entities.Example, error) {
	if len(id) == 0 {
		return entities.Example{}, fmt.Errorf("missing id: %w", domain.ErrMalformedParameters)
	}

	example, err := uc.R.GetExampleByID(ctx, id)
	if err != nil {
		return entities.Example{}, fmt.Errorf("failed to get example by id: %w", err)
	}

	return example, nil
}
