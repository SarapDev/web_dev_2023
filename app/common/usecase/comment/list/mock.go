package list

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

var _ CommentListExecutor = MockExecutor{}

type (
	MockExecutor struct {
		MockExecuteFunc MockExecuteFunc
	}
	MockExecuteFunc func(
		ctx context.Context,
		postID int,
	) ([]*entity.Comment, error)
)

func (m MockExecutor) Execute(
	ctx context.Context,
	postID int,
) ([]*entity.Comment, error) {
	return m.MockExecuteFunc(ctx, postID)
}
