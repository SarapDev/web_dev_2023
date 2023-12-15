package get

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"go.uber.org/multierr"
	"golang.org/x/exp/slog"
)

var (
	ErrPostNotFound  = errors.New("пост не найден")
	ErrGetPostFailed = errors.New("ошибка при получения поста")
)

type PostGetExecutor interface {
	Execute(ctx context.Context, id int) (*entity.Post, error)
}

var _ PostGetExecutor = Executor{}

type Executor struct {
	repo   post.Repository
	tracer opentracing.Tracer
	logger *slog.Logger
}

func NewExecutor(
	repo post.Repository,
	tracer opentracing.Tracer,
	logger *slog.Logger,
) *Executor {
	return &Executor{
		repo:   repo,
		tracer: tracer,
		logger: logger,
	}
}

func (e Executor) Execute(ctx context.Context, id int) (*entity.Post, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, e.tracer,
		"app::common::usecase::post::get::executor::execute",
	)
	defer span.Finish()

	e.logger.Info("выполнение получения поста", "id", id)

	span.SetTag("id", id)

	object, err := e.repo.Get(ctx, id)

	switch {
	case errors.Is(err, repository.ErrNotFound):
		return nil, ErrPostNotFound
	case err != nil:
		return nil, multierr.Combine(ErrGetPostFailed, err)
	}

	e.logger.Info("пост успешно получен", "id", id)

	return object, nil
}
