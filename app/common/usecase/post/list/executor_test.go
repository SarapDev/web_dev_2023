package list

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"testing"
	"time"
)

func TestExecutor_Execute(t *testing.T) {
	tests := []struct {
		name string

		repo   post.Repository
		tracer opentracing.Tracer

		offset int
		limit  int

		want func(t *testing.T, got []*entity.Post, total int, err error)
	}{
		{
			name: "успешное получение списка постов",

			repo: post.MockRepository{
				MockListFunc: func(
					ctx context.Context,
					offset int,
					limit int,
				) ([]*entity.Post, error) {
					assert.Equal(t, 15, offset)
					assert.Equal(t, 1, limit)

					return []*entity.Post{
						{
							ID:        1,
							Title:     "title",
							Content:   "content",
							CreatedAt: time.UnixMicro(1028301982),
						},
					}, nil
				},
				MockCountFunc: func(
					ctx context.Context,
				) (int, error) {
					return 146, nil
				},
			},

			offset: 15,
			limit:  1,

			want: func(t *testing.T, got []*entity.Post, total int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 1, len(got))
				assert.Equal(t, 146, total)
				assert.Equal(t, 1, got[0].ID)
				assert.Equal(t, "title", got[0].Title)
				assert.Equal(t, "content", got[0].Content)
				assert.Equal(t, time.UnixMicro(1028301982), got[0].CreatedAt)
			},
		},
		{
			name: "ошибка получения списка постов",

			repo: post.MockRepository{
				MockListFunc: func(
					ctx context.Context,
					offset int,
					limit int,
				) ([]*entity.Post, error) {
					return nil, assert.AnError
				},
			},

			offset: 15,
			limit:  1,

			want: func(t *testing.T, got []*entity.Post, total int, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, ErrGetPostsFailed)
			},
		},
		{
			name: "ошибка получения количества постов",

			repo: post.MockRepository{
				MockListFunc: func(
					ctx context.Context,
					offset int,
					limit int,
				) ([]*entity.Post, error) {
					return []*entity.Post{
						{
							ID:        1,
							Title:     "title",
							Content:   "content",
							CreatedAt: time.UnixMicro(1028301982),
						},
					}, nil
				},
				MockCountFunc: func(
					ctx context.Context,
				) (int, error) {
					return 0, assert.AnError
				},
			},

			offset: 15,
			limit:  1,

			want: func(t *testing.T, got []*entity.Post, total int, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, ErrGetPostsTotalFailed)
			},
		},
		{
			name: "ошибка что нет постов",

			repo: post.MockRepository{
				MockListFunc: func(
					ctx context.Context,
					offset int,
					limit int,
				) ([]*entity.Post, error) {
					return nil, nil
				},
			},

			offset: 15,
			limit:  1,

			want: func(t *testing.T, got []*entity.Post, total int, err error) {
				assert.ErrorIs(t, err, ErrPostsNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var (
					executor         = NewExecutor(tt.repo, tt.tracer)
					list, total, err = executor.Execute(
						context.Background(),
						tt.offset,
						tt.limit,
					)
				)

				tt.want(t, list, total, err)
			},
		)
	}
}
