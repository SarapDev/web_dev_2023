package list

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"github.kolesa-team.org/backend/go-module/logger"
	"go.uber.org/multierr"
)

var (
	ErrGetPostsFailed      = errors.New("ошибка получения списка постов")
	ErrPostsNotFound       = errors.New("посты не найдены")
	ErrGetPostsTotalFailed = errors.New("ошибка получения количества постов")
)

type PostListExecutor interface {
	Execute(
		ctx context.Context,
		offset, limit int,
	) ([]*entity.Post, int, error)
}

var _ PostListExecutor = Executor{}

type Executor struct {
	repository post.Repository
	tracer     opentracing.Tracer
}

func NewExecutor(
	repository post.Repository,
	tracer opentracing.Tracer,
) *Executor {
	return &Executor{repository: repository, tracer: tracer}
}

func (e Executor) Execute(
	ctx context.Context,
	offset, limit int,
) (list []*entity.Post, total int, err error) {
	span := opentracing.StartSpan(
		"app::common::usecase::post::list::executor::execute",
	)
	defer span.Finish()

	logger.FromContextOrDiscard(ctx).
		Info(
			"получение списка постов",
			"offset", offset,
			"limit", limit,
		)

	span.LogFields(
		log.Int("offset", offset),
		log.Int("limit", limit),
	)

	objects, err := e.repository.List(ctx, offset, limit)

	switch {
	case err != nil:
		return nil, 0, multierr.Append(ErrGetPostsFailed, err)
	case len(objects) == 0:
		return nil, 0, ErrPostsNotFound
	}

	total, err = e.repository.Count(ctx)

	if err != nil {
		return nil, 0, multierr.Append(ErrGetPostsTotalFailed, err)
	}

	return objects, total, err
}
