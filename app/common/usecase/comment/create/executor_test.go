package create

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/comment"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/get"
	"golang.org/x/exp/slog"
	"io"
	"testing"
	"time"
)

func TestExecutor_Execute(t *testing.T) {
	tests := []struct {
		name string

		repo            comment.Repository
		getPostExecutor get.PostGetExecutor

		postID  int
		author  string
		content string

		want func(t *testing.T, got *entity.Comment, err error)
	}{
		{
			name: "успешное создание комментария",

			repo: comment.MockRepository{
				MockCreateFunc: func(
					ctx context.Context,
					comment *entity.Comment,
				) error {
					assert.Equal(t, 1, comment.PostID)
					assert.Equal(t, "author", comment.Author)
					assert.Equal(t, "content", comment.Content)
					assert.GreaterOrEqual(
						t,
						time.Now().Unix(),
						comment.CreatedAt.Unix(),
					)

					comment.ID = 2

					return nil
				},
			},
			getPostExecutor: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					postID int,
				) (*entity.Post, error) {
					assert.Equal(t, 1, postID)

					return &entity.Post{
						ID: postID,
					}, nil
				},
			},

			postID:  1,
			author:  "author",
			content: "content",

			want: func(t *testing.T, got *entity.Comment, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 2, got.ID)
				assert.Equal(t, 1, got.PostID)
				assert.Equal(t, "author", got.Author)
				assert.Equal(t, "content", got.Content)
				assert.GreaterOrEqual(
					t,
					time.Now().Unix(),
					got.CreatedAt.Unix(),
				)
			},
		},
		{
			name: "пост не найден",

			getPostExecutor: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					postID int,
				) (*entity.Post, error) {
					return nil, get.ErrPostNotFound
				},
			},

			postID: 1,

			want: func(t *testing.T, got *entity.Comment, err error) {
				assert.ErrorIs(t, err, ErrPostNotFound)
			},
		},
		{
			name: "ошибка при создании комментария",

			repo: comment.MockRepository{
				MockCreateFunc: func(
					ctx context.Context,
					comment *entity.Comment,
				) error {
					return assert.AnError
				},
			},
			getPostExecutor: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					postID int,
				) (*entity.Post, error) {
					return &entity.Post{
						ID: postID,
					}, nil
				},
			},

			postID:  1,
			author:  "author",
			content: "content",

			want: func(t *testing.T, got *entity.Comment, err error) {
				assert.ErrorIs(t, err, assert.AnError)
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
						tt.getPostExecutor,
					)
					got, err = executor.Execute(
						context.Background(),
						tt.postID,
						tt.author,
						tt.content,
					)
				)

				tt.want(t, got, err)
			},
		)
	}
}
