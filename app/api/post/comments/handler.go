package comments

import (
	"errors"
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/comment/list"
	chicustom "github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

const (
	ErrCodeInvalidPostID     = 1001
	ErrCodePostNotFound      = 1002
	ErrCodeGetCommentsFailed = 1003
)

type entry struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	usecase list.CommentListExecutor
	tracer  opentracing.Tracer
}

func NewHandler(
	usecase list.CommentListExecutor,
	tracer opentracing.Tracer,
) Handler {
	return Handler{
		usecase: usecase,
		tracer:  tracer,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span, logger := helper.GetRequestTools(
		h.tracer, r, moduleName,
	)
	defer span.Finish()

	id := chicustom.GetIntFromRequest(r, "post_id", 0)
	if id == 0 {
		response.SetJSONResponse(
			w, http.StatusBadRequest, response.WithError(
				ErrCodeInvalidPostID,
				"не валидный id поста",
				"не валидный id поста",
			),
		)

		return
	}

	span.SetTag("post_id", id)

	comments, err := h.usecase.Execute(ctx, id)

	switch {
	case errors.Is(err, list.ErrPostNotFound):
		ext.LogError(span, err)
		logger.Error(
			"пост не найден",
			"error", err,
			"id", id,
		)

		response.SetJSONResponse(
			w, http.StatusNotFound, response.WithError(
				ErrCodePostNotFound,
				"пост не найден",
				err.Error(),
			),
		)

		return
	case err != nil:
		ext.LogError(span, err)
		logger.Error(
			"ошибка получения поста",
			"error", err,
			"id", id,
		)

		response.SetJSONResponse(
			w, http.StatusInternalServerError, response.WithError(
				ErrCodeGetCommentsFailed,
				"ошибка получения поста",
				err.Error(),
			),
		)

		return
	}

	logger.Info(
		"получены комментарии",
		"post_id", id,
		"count", len(comments),
	)

	data := make([]entry, 0, len(comments))

	for _, c := range comments {
		data = append(
			data, entry{
				ID:        c.ID,
				Content:   c.Content,
				Author:    c.Author,
				CreatedAt: c.CreatedAt,
			},
		)
	}

	response.SetJSONResponse(
		w, http.StatusOK,
		response.WithData(data),
	)
}
