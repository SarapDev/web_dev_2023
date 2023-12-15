package get

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"golang.org/x/exp/slog"
	"io"
	"testing"
	"time"
)

func TestExecutor_Execute(t *testing.T) {
	tests := []struct {
		name string

		repo   post.Repository
		postID int

		want func(t *testing.T, got *entity.Post, err error)
	}{
		{
			name: "успешное получение поста",

			repo: post.MockRepository{
				MockGetFunc: func(
					ctx context.Context,
					postID int,
				) (*entity.Post, error) {
					assert.Equal(t, 1, postID)

					return &entity.Post{
						ID:        postID,
						Title:     "title",
						Content:   "content",
						CreatedAt: time.UnixMicro(1028301982),
					}, nil
				},
			},
			postID: 1,

			want: func(t *testing.T, got *entity.Post, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, got.ID)
				assert.Equal(t, "title", got.Title)
				assert.Equal(t, "content", got.Content)
				assert.Equal(t, time.UnixMicro(1028301982), got.CreatedAt)
			},
		},
		{
			name: "ошибка получения поста",

			repo: post.MockRepository{
				MockGetFunc: func(
					ctx context.Context,
					postID int,
				) (*entity.Post, error) {
					return nil, assert.AnError
				},
			},
			postID: 1,

			want: func(t *testing.T, got *entity.Post, err error) {
				assert.Error(t, err)
				assert.Nil(t, got)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var (
					executor = NewExecutor(
						tt.repo,
						opentracing.NoopTracer{},
						slog.New(slog.NewJSONHandler(io.Discard, nil)),
					)
					got, err = executor.Execute(
						context.Background(),
						tt.postID,
					)
				)

				tt.want(t, got, err)
			},
		)
	}
}
