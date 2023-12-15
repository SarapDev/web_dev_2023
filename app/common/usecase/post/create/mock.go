package create

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

var _ PostCreateExecutor = MockExecutor{}

type (
	MockExecutor struct {
		MockExecuteFunc MockExecuteFunc
	}
	MockExecuteFunc func(
		ctx context.Context,
		title string,
		content string,
	) (*entity.Post, error)
)

func (m MockExecutor) Execute(
	ctx context.Context,
	title string,
	content string,
) (*entity.Post, error) {
	return m.MockExecuteFunc(ctx, title, content)
}
