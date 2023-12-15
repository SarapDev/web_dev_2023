package external

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/wagslane/go-rabbitmq"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/create"
	"golang.org/x/exp/slog"
)

type req struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Handler struct {
	logger  *slog.Logger
	usecase create.PostCreateExecutor
	tracer  opentracing.Tracer
}

func NewHandler(
	logger *slog.Logger,
	usecase create.PostCreateExecutor,
	tracer opentracing.Tracer,
) Handler {
	return Handler{
		logger:  logger,
		usecase: usecase,
		tracer:  tracer,
	}
}

func (h Handler) Handle(
	ctx context.Context,
	delivery rabbitmq.Delivery,
) rabbitmq.Action {
	span := h.tracer.StartSpan(
		"app::external::handler",
		opentracing.ChildOf(
			opentracing.SpanFromContext(ctx).Context(),
		),
	)
	defer span.Finish()

	h.logger.Info(
		"обработка сообщения из очереди",
		"body", string(delivery.Body),
	)

	var r req

	if err := json.Unmarshal(delivery.Body, &r); err != nil {
		ext.LogError(span, err)
		h.logger.Error(
			"не удалось распарсить тело сообщения",
			"error", err,
			"body", string(delivery.Body),
		)

		return rabbitmq.NackDiscard
	}

	post, err := h.usecase.Execute(
		ctx, r.Title, r.Content,
	)

	if err != nil {
		ext.LogError(span, err)
		h.logger.Error(
			"не удалось создать пост",
			"error", err,
			"request", r,
		)

		return rabbitmq.NackRequeue
	}

	h.logger.Info(
		"пост успешно создан через external worker",
		"post", *post,
	)

	return rabbitmq.Ack
}
