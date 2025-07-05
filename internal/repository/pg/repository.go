package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-template/domain"
	"go-template/domain/entities"
	"go-template/internal/repository/pg/gen"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository implements a domain.ExampleRepository interface.
type Repository struct {
	queries *gen.Queries
	db      *pgxpool.Pool
}

// NewRepository creates a new Repository instance.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		queries: gen.New(db),
		db:      db,
	}
}

// CreateExample creates a new example in the database.
func (r *Repository) CreateExample(ctx context.Context, input entities.Example) (string, error) {
	out, err := r.queries.CreateExample(ctx, input.Title, input.Content)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("example with title '%s' already exists: %w", input.Title, domain.ErrDuplicateKey)
		}
		return "", err
	}

	return out.String(), nil
}

// GetExampleByID retrieves an example by its ID.
func (r *Repository) GetExampleByID(ctx context.Context, id string) (entities.Example, error) {
	out, err := r.queries.GetExampleByID(ctx, uuid.FromStringOrNil(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Example{}, nil
		}
		return entities.Example{}, err
	}

	return entities.Example{
		ID:        out.ID.String(),
		Title:     out.Title,
		Content:   out.Content,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}, nil
}
