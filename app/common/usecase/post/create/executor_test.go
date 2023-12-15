package create

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"golang.org/x/exp/slog"
	"io"
	"testing"
)

func TestExecutor_Execute(t *testing.T) {
	tests := []struct {
		name string
		repo post.Repository

		argTitle   string
		argContent string

		want func(t *testing.T, got *entity.Post, err error)
	}{
		{
			name: "успешное создание поста",
			repo: post.MockRepository{
				MockCreateFunc: func(ctx context.Context, post *entity.Post) error {
					post.ID = 1

					assert.Equal(t, "hello world", post.Title)
					assert.Equal(t, "hi my friend. how are you?", post.Content)

					return nil
				},
			},
			argTitle:   "hello world",
			argContent: "hi my friend. how are you?",

			want: func(t *testing.T, got *entity.Post, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, got.ID)
			},
		},
		{
			name: "ошибка при создании поста",
			repo: post.MockRepository{
				MockCreateFunc: func(ctx context.Context, post *entity.Post) error {
					return assert.AnError
				},
			},
			argTitle:   "hello world",
			argContent: "hi my friend. how are you?",

			want: func(t *testing.T, got *entity.Post, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				executor := NewExecutor(
					opentracing.NoopTracer{},
					slog.New(slog.NewJSONHandler(io.Discard, nil)),
					tt.repo,
				)

				got, err := executor.Execute(
					context.Background(),
					tt.argTitle,
					tt.argContent,
				)

				tt.want(t, got, err)
			},
		)
	}
}
