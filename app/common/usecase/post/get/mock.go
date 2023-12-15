package get

import (
	"context"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
)

var _ PostGetExecutor = MockExecutor{}

type (
	MockExecutor struct {
		MockExecuteFunc MockExecuteFunc
	}
	MockExecuteFunc func(ctx context.Context, id int) (*entity.Post, error)
)

func (m MockExecutor) Execute(ctx context.Context, id int) (*entity.Post, error) {
	return m.MockExecuteFunc(ctx, id)
}
