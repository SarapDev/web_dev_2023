package list

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

var _ PostListExecutor = MockExecutor{}

type (
	MockExecutor struct {
		MockExecuteFunc MockExecuteFunc
	}
	MockExecuteFunc func(
		ctx context.Context,
		offset, limit int,
	) ([]*entity.Post, int, error)
)

func (m MockExecutor) Execute(
	ctx context.Context,
	offset, limit int,
) ([]*entity.Post, int, error) {
	return m.MockExecuteFunc(ctx, offset, limit)
}
