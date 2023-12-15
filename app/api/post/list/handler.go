package list

import (
	"errors"
	"github.kolesa-team.org/backend/go-module/chi/helper"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/samber/lo"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/list"
	chicustom "github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/chi/response"
)

const (
	maxLimit     = 100
	maxOffset    = 1000
	defaultLimit = 10
)

const (
	ErrCodePostsNotFound  = 1001
	ErrCodeGetPostsFailed = 1002
)

type entry struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	usecase list.PostListExecutor
	tracer  opentracing.Tracer
}

func NewHandler(
	usecase list.PostListExecutor,
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

	offset := chicustom.GetIntFromRequest(r, "offset", 0)
	limit := chicustom.GetIntFromRequest(r, "offset", maxLimit)

	// ограничиваем 0 >= offset <= maxOffset
	offset = lo.Max([]int{lo.Min([]int{offset, maxOffset}), 0})
	// ограничиваем defaultLimit >= limit <= maxLimit
	limit = lo.Max([]int{lo.Min([]int{limit, maxLimit}), defaultLimit})

	span.LogFields(
		log.Int("offset", offset),
		log.Int("limit", limit),
	)

	posts, total, err := h.usecase.Execute(ctx, offset, limit)

	switch {
	case errors.Is(err, list.ErrPostsNotFound):
		logger.Error(
			"посты не найдены",
			"error", err,
			"offset", offset,
			"limit", limit,
		)

		response.SetJSONResponse(
			w, http.StatusNotFound, response.WithError(
				ErrCodePostsNotFound,
				"посты не найдены",
				err.Error(),
			),
		)

		return
	case err != nil:
		logger.Error(
			"ошибка получения постов",
			"error", err,
			"offset", offset,
			"limit", limit,
		)

		response.SetJSONResponse(
			w, http.StatusInternalServerError, response.WithError(
				ErrCodeGetPostsFailed,
				"ошибка получения постов",
				err.Error(),
			),
		)

		return
	}

	span.LogFields(
		log.Int("total", total),
	)

	logger.Info(
		"посты получены",
		"offset", offset,
		"limit", limit,
		"total", total,
	)

	data := make([]entry, 0, len(posts))

	for _, p := range posts {
		data = append(
			data, entry{
				ID:        p.ID,
				Title:     p.Title,
				CreatedAt: p.CreatedAt,
			},
		)
	}

	response.SetJSONResponse(
		w, http.StatusOK,
		response.WithData(data),
		response.WithMeta(
			map[string]any{
				"total": total,
			},
		),
	)
}
