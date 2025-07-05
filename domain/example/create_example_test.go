package example

import (
	"context"
	"testing"

	"finance/domain/entities"
	"finance/domain/example/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateExample(t *testing.T) {
	tests := []struct {
		name    string
		input   entities.Example
		mock    func(*mocks.RepositoryMock)
		wantErr bool
	}{
		{
			name: "success",
			input: entities.Example{
				ID:    "123",
				Title: "Test Title",
			},
			mock: func(m *mocks.RepositoryMock) {
				m.CreateExampleFunc = func(ctx context.Context, input entities.Example) (string, error) {
					return "123", nil
				}
			},
			wantErr: false,
		},
		{
			name: "empty title",
			input: entities.Example{
				ID: "123",
			},
			mock:    func(m *mocks.RepositoryMock) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			tt.mock(repo)

			uc := New(repo)
			id, err := uc.CreateExample(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.ID, id)
			}
		})
	}
}
