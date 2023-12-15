package list

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

		postID int

		want func(t *testing.T, got []*entity.Comment, err error)
	}{
		{
			name: "успешное получение списка комментариев",

			repo: comment.MockRepository{
				MockListByPostIDFunc: func(
					ctx context.Context,
					postID int,
				) ([]*entity.Comment, error) {
					assert.Equal(t, 1, postID)

					return []*entity.Comment{
						{
							ID:        1,
							PostID:    1,
							Author:    "bob",
							Content:   "bar",
							CreatedAt: time.Now(),
						},
						{
							ID:      2,
							PostID:  1,
							Author:  "alica",
							Content: "foo",
						},
					}, nil
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
			postID: 1,

			want: func(t *testing.T, got []*entity.Comment, err error) {
				assert.NoError(t, err)
				assert.Len(t, got, 2)

				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, 1, got[0].PostID)
				assert.Equal(t, "bob", got[0].Author)
				assert.Equal(t, "bar", got[0].Content)
				assert.GreaterOrEqual(
					t,
					time.Now().Unix(),
					got[0].CreatedAt.Unix(),
				)

				assert.Equal(t, 2, got[1].ID)
				assert.Equal(t, 1, got[1].PostID)
				assert.Equal(t, "alica", got[1].Author)
				assert.Equal(t, "foo", got[1].Content)
				assert.Equal(t, time.Time{}, got[1].CreatedAt)
			},
		},
		{
			name: "пост не найден",

			getPostExecutor: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					postID int,
				) (*entity.Post, error) {
					assert.Equal(t, 1, postID)

					return nil, get.ErrPostNotFound
				},
			},
			postID: 1,

			want: func(t *testing.T, got []*entity.Comment, err error) {
				assert.ErrorIs(t, err, ErrPostNotFound)
				assert.Nil(t, got)
			},
		},
		{
			name: "ошибка при получении поста",

			repo: comment.MockRepository{
				MockListByPostIDFunc: func(
					ctx context.Context,
					postID int,
				) ([]*entity.Comment, error) {
					assert.Equal(t, 1, postID)

					return nil, io.EOF
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

			postID: 1,

			want: func(t *testing.T, got []*entity.Comment, err error) {
				assert.ErrorIs(t, err, io.EOF)
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
						tt.getPostExecutor,
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
