package create

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/comment"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/get"
	"go.uber.org/multierr"
	"golang.org/x/exp/slog"
	"time"
)

var (
	ErrCreateCommentFailed = errors.New("ошибка создания комментария")
	ErrPostNotFound        = errors.New("пост не найден")
)

type CommentCreateExecutor interface {
	Execute(
		ctx context.Context,
		postID int,
		author string,
		content string,
	) (*entity.Comment, error)
}

var _ CommentCreateExecutor = Executor{}

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
	author string,
	content string,
) (*entity.Comment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, e.tracer,
		"app::common::usecase::comment::create::executor::execute",
	)
	defer span.Finish()

	e.logger.Info(
		"выполнение создание комментария",
		"postID", postID,
		"author", author,
		"content", content,
	)

	span.SetTag("post_id", postID)
	span.LogFields(
		log.String("author", author),
		log.String("content", content),
	)

	_, err := e.getPostExecutor.Execute(ctx, postID)

	switch {
	case errors.Is(err, get.ErrPostNotFound):
		return nil, multierr.Append(ErrPostNotFound, err)
	case err != nil:
		return nil, multierr.Append(ErrCreateCommentFailed, err)
	}

	entry := &entity.Comment{
		PostID:    postID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err = e.repo.Create(ctx, entry); err != nil {
		return nil, multierr.Append(ErrCreateCommentFailed, err)
	}

	e.logger.Info(
		"комментарий успешно создан",
		"commentID", entry.ID,
	)

	return entry, nil
}
