package comment

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

type Repository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	ListByPostID(ctx context.Context, postID int) ([]*entity.Comment, error)
}
