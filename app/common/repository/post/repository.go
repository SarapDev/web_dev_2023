package post

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

type Repository interface {
	Create(ctx context.Context, post *entity.Post) error
	Get(ctx context.Context, id int) (*entity.Post, error)
	List(ctx context.Context, offset, limit int) ([]*entity.Post, error)
	Count(ctx context.Context) (int, error)
}
