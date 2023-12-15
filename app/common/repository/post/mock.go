package post

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

var _ Repository = MockRepository{}

type (
	MockCreateFunc func(
		ctx context.Context,
		post *entity.Post,
	) error
	MockGetFunc func(
		ctx context.Context,
		id int,
	) (*entity.Post, error)
	MockListFunc func(
		ctx context.Context,
		offset, limit int,
	) ([]*entity.Post, error)
	MockCountFunc func(
		ctx context.Context,
	) (int, error)
)

type MockRepository struct {
	MockCreateFunc MockCreateFunc
	MockGetFunc    MockGetFunc
	MockListFunc   MockListFunc
	MockCountFunc  MockCountFunc
}

func (m MockRepository) Create(ctx context.Context, post *entity.Post) error {
	return m.MockCreateFunc(ctx, post)
}

func (m MockRepository) Get(ctx context.Context, id int) (*entity.Post, error) {
	return m.MockGetFunc(ctx, id)
}

func (m MockRepository) List(ctx context.Context, offset, limit int) ([]*entity.Post, error) {
	return m.MockListFunc(ctx, offset, limit)
}

func (m MockRepository) Count(ctx context.Context) (int, error) {
	return m.MockCountFunc(ctx)
}
