package comment

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

type (
	MockCreateFunc func(
		ctx context.Context,
		comment *entity.Comment,
	) error
	MockListByPostIDFunc func(
		ctx context.Context,
		postID int,
	) ([]*entity.Comment, error)
)

type MockRepository struct {
	MockCreateFunc       MockCreateFunc
	MockListByPostIDFunc MockListByPostIDFunc
}

func (m MockRepository) Create(
	ctx context.Context,
	comment *entity.Comment,
) error {
	return m.MockCreateFunc(ctx, comment)
}

func (m MockRepository) ListByPostID(
	ctx context.Context,
	postID int,
) ([]*entity.Comment, error) {
	return m.MockListByPostIDFunc(ctx, postID)
}
