package create

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

var _ CommentCreateExecutor = MockExecutor{}

type (
	MockExecutor struct {
		MockExecuteFunc MockExecuteFunc
	}
	MockExecuteFunc func(
		ctx context.Context,
		postID int,
		author string,
		content string,
	) (*entity.Comment, error)
)

func (m MockExecutor) Execute(
	ctx context.Context,
	postID int,
	author string,
	content string,
) (*entity.Comment, error) {
	return m.MockExecuteFunc(ctx, postID, author, content)
}
