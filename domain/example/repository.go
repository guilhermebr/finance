package example

import (
	"context"
	"go-template/domain/entities"
)

//go:generate moq -skip-ensure -stub -pkg mocks -out mocks/repository.go . Repository
type Repository interface {
	CreateExample(context.Context, entities.Example) (string, error)
	GetExampleByID(context.Context, string) (entities.Example, error)
}
