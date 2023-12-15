package create

import (
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/create"
	chicustom "github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

const (
	ErrCodeInvalidBody      = 1001
	ErrCodeCannotReadBody   = 1002
	ErrCodeCreatePostFailed = 1003
)

type entry struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Handler struct {
	usecase create.PostCreateExecutor
	tracer  opentracing.Tracer
}

func NewHandler(
	usecase create.PostCreateExecutor,
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

	data, err := chicustom.ParseRequestJSON[Data](r)
	if err != nil {
		response.SetJSONResponse(
			w, http.StatusBadRequest, response.WithError(
				ErrCodeCannotReadBody,
				"invalid request body",
				err.Error(),
			),
		)

		return
	}

	if errs := data.Validate(); len(errs) > 0 {
		ext.LogError(span, err)
		logger.Warn(
			"ошибка при валидации тела запроса",
			"errors", errs,
			"data", data,
		)
		response.SetJSONResponse(
			w, http.StatusBadRequest,
			response.WithErrors(errs),
		)

		return
	}

	post, err := h.usecase.Execute(ctx, data.Title, data.Content)

	if err != nil {
		ext.LogError(span, err)
		logger.Error(
			"ошибка при создании поста",
			"error", err,
		)

		response.SetJSONResponse(
			w, http.StatusInternalServerError, response.WithError(
				ErrCodeCreatePostFailed,
				"internal server error",
				err.Error(),
			),
		)

		return
	}

	span.SetTag("post_id", post.ID)

	logger.Info(
		"пост успешно создан",
		"id", post.ID,
	)

	response.SetJSONResponse(
		w, http.StatusCreated,
		response.WithData(
			entry{
				ID:      post.ID,
				Title:   post.Title,
				Content: post.Content,
			},
		),
	)
}
