package create

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"go.uber.org/multierr"
	"golang.org/x/exp/slog"
	"time"
)

var ErrCreatePostFailed = errors.New("ошибка создания поста")

type PostCreateExecutor interface {
	Execute(
		ctx context.Context,
		title string,
		content string,
	) (*entity.Post, error)
}

var _ PostCreateExecutor = Executor{}

type Executor struct {
	tracer opentracing.Tracer
	logger *slog.Logger
	repo   post.Repository
}

func NewExecutor(
	tracer opentracing.Tracer,
	logger *slog.Logger,
	repo post.Repository,
) *Executor {
	return &Executor{
		tracer: tracer,
		logger: logger,
		repo:   repo,
	}
}

func (e Executor) Execute(
	ctx context.Context,
	title string,
	content string,
) (*entity.Post, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, e.tracer,
		"app::common::usecase::post::create::executor::execute",
	)
	defer span.Finish()

	e.logger.Info(
		"выполнение создание поста",
		"title", title,
		"content", content,
	)

	span.LogFields(
		log.String("title", title),
		log.String("content", content),
	)

	entry := &entity.Post{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := e.repo.Create(ctx, entry); err != nil {
		return nil, multierr.Append(ErrCreatePostFailed, err)
	}

	span.SetTag("post_id", entry.ID)

	e.logger.Info(
		"пост успешно создан",
		"post_id", entry.ID,
	)

	return entry, nil
}
