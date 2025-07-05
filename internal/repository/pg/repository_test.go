package pg

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go-template/domain/entities"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=go_app_template_test",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {{HostIP: "localhost", HostPort: "5432"}},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		panic(fmt.Sprintf("Could not start resource: %s", err))
	}

	// Wait for the container to be ready
	if err := pool.Retry(func() error {
		conn, err := pgxpool.New(context.Background(), getTestDSN())
		if err != nil {
			return err
		}
		return conn.Ping(context.Background())
	}); err != nil {
		panic(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	// Run migrations
	mig, err := migrate.New(
		"file://migrations",
		getTestDSN(),
	)
	if err != nil {
		panic(fmt.Sprintf("Could not create migration: %s", err))
	}

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("Could not run migrations: %s", err))
	}

	code := m.Run()

	// Cleanup
	if err := pool.Purge(resource); err != nil {
		panic(fmt.Sprintf("Could not purge resource: %s", err))
	}

	os.Exit(code)
}

func getTestDSN() string {
	return fmt.Sprintf("postgres://postgres:postgres@localhost:%s/go_app_template_test?sslmode=disable", resource.GetPort("5432/tcp"))
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), getTestDSN())
	require.NoError(t, err)
	return pool
}

func TestRepository_CreateExample(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewRepository(pool)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   entities.Example
		wantErr bool
	}{
		{
			name: "success",
			input: entities.Example{
				ID:    uuid.Must(uuid.NewV4()).String(),
				Title: "Test Title",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.CreateExample(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the record was created
				got, err := repo.GetExampleByID(ctx, id)
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Title, got.Title)
			}
		})
	}
}
