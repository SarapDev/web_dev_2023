package get

import (
	"errors"
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/get"
	chicustom "github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

const (
	ErrCodeInvalidPostID = 1001
	ErrCodePostNotFound  = 1002
	ErrCodeGetPostFailed = 1003
)

type entry struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	usecase get.PostGetExecutor
	tracer  opentracing.Tracer
}

func NewHandler(
	usecase get.PostGetExecutor,
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

	post, err := h.usecase.Execute(ctx, id)

	switch {
	case errors.Is(err, get.ErrPostNotFound):
		ext.LogError(span, err)
		logger.Warn(
			"пост не найден",
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
			"ошибка при получении поста",
			"error", err,
			"id", id,
		)

		response.SetJSONResponse(
			w, http.StatusInternalServerError, response.WithError(
				ErrCodeGetPostFailed,
				"ошибка при получении поста",
				err.Error(),
			),
		)

		return
	}

	logger.Info(
		"пост получен",
		"id", post.ID,
	)

	response.SetJSONResponse(
		w, http.StatusOK, response.WithData(
			entry{
				ID:        post.ID,
				Title:     post.Title,
				Content:   post.Content,
				CreatedAt: post.CreatedAt,
			},
		),
	)
}
