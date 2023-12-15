package list

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/comment"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/get"
	"go.uber.org/multierr"
	"golang.org/x/exp/slog"
)

var (
	ErrPostNotFound      = errors.New("пост не найден")
	ErrGetCommentsFailed = errors.New("ошибка получения списка комментариев")
)

type CommentListExecutor interface {
	Execute(ctx context.Context, postID int) ([]*entity.Comment, error)
}

type Executor struct {
	repo            comment.Repository
	tracer          opentracing.Tracer
	logger          *slog.Logger
	getPostExecutor get.PostGetExecutor
}

func NewExecutor(
	repo comment.Repository,
	tracer opentracing.Tracer,
	logger *slog.Logger,
	getPostExecutor get.PostGetExecutor,
) *Executor {
	return &Executor{
		repo:            repo,
		tracer:          tracer,
		logger:          logger,
		getPostExecutor: getPostExecutor,
	}
}

func (e Executor) Execute(
	ctx context.Context,
	postID int,
) ([]*entity.Comment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, e.tracer,
		"app::common::usecase::comment::list::executor::execute",
	)
	defer span.Finish()

	e.logger.Info(
		"получение списка комментариев",
		"postID", postID,
	)

	span.SetTag("post_id", postID)

	_, err := e.getPostExecutor.Execute(ctx, postID)

	switch {
	case errors.Is(err, get.ErrPostNotFound):
		return nil, multierr.Combine(ErrPostNotFound, err)
	case err != nil:
		return nil, multierr.Combine(ErrGetCommentsFailed, err)
	}

	comments, err := e.repo.ListByPostID(ctx, postID)

	if err != nil {
		return nil, multierr.Combine(ErrGetCommentsFailed, err)
	}

	e.logger.Info(
		"список комментариев успешно получен",
		"postID", postID,
		"count", len(comments),
	)

	return comments, nil
}
