package example

import (
	"context"
	"testing"

	"go-template/domain/entities"
	"go-template/domain/example/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetExampleByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func(*mocks.RepositoryMock)
		want    entities.Example
		wantErr bool
	}{
		{
			name: "success",
			id:   "123",
			mock: func(m *mocks.RepositoryMock) {
				m.GetExampleByIDFunc = func(ctx context.Context, id string) (entities.Example, error) {
					return entities.Example{
						ID:    "123",
						Title: "Test Title",
					}, nil
				}
			},
			want: entities.Example{
				ID:    "123",
				Title: "Test Title",
			},
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			mock:    func(m *mocks.RepositoryMock) {},
			want:    entities.Example{},
			wantErr: true,
		},
		{
			name: "not found",
			id:   "999",
			mock: func(m *mocks.RepositoryMock) {
				m.GetExampleByIDFunc = func(ctx context.Context, id string) (entities.Example, error) {
					return entities.Example{}, nil
				}
			},
			want:    entities.Example{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			tt.mock(repo)

			uc := New(repo)
			got, err := uc.GetExampleByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
