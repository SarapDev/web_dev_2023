package create

import (
	"errors"
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/comment/create"
	chicustom "github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

const (
	ErrCodeCannotReadBody      = 1001
	ErrCodeInvalidBody         = 1002
	ErrCodePostNotFound        = 1003
	ErrCodeCreateCommentFailed = 1004
)

type entry struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type Handler struct {
	usecase create.CommentCreateExecutor
	tracer  opentracing.Tracer
}

func NewHandler(
	usecase create.CommentCreateExecutor,
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
				"ошибка при чтении тела запроса",
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

	comment, err := h.usecase.Execute(ctx, data.PostID, data.Author, data.Content)

	switch {
	case errors.Is(err, create.ErrPostNotFound):
		ext.LogError(span, err)
		logger.Error(
			"пост не найден",
			"error", err,
			"post_id", data.PostID,
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
			"ошибка при создании комментария",
			"error", err,
		)

		response.SetJSONResponse(
			w, http.StatusInternalServerError, response.WithError(
				ErrCodeCreateCommentFailed,
				"ошибка при создании комментария",
				err.Error(),
			),
		)

		return
	}

	span.SetTag("comment_id", comment.ID)

	logger.Info(
		"комментарий успешно создан",
		"comment_id", comment.ID,
	)

	response.SetJSONResponse(
		w, http.StatusCreated,
		response.WithData(
			entry{
				ID:      comment.ID,
				PostID:  comment.PostID,
				Author:  comment.Author,
				Content: comment.Content,
			},
		),
	)
}
