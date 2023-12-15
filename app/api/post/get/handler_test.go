package get

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/get"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name    string
		usecase get.PostGetExecutor
		request func() *http.Request
		want    func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "успешное получение поста",
			usecase: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					id int,
				) (*entity.Post, error) {
					assert.Equal(t, 1, id)

					return &entity.Post{
						ID:        1,
						Title:     "title",
						Content:   "content",
						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/post/1",
					nil,
				)
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("post_id", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.JSONEq(
					t,
					`{"data": {"id":1,"title":"title","content":"content", "created_at":"2021-01-01T00:00:00Z"}}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "пост не найден",

			usecase: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					id int,
				) (*entity.Post, error) {
					assert.Equal(t, 1, id)

					return nil, get.ErrPostNotFound
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/post/1",
					nil,
				)
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("post_id", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, resp.Code)
				assert.JSONEq(
					t,
					`{"errors": [{"status":1002, "title": "пост не найден", "detail": "пост не найден"}]}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "ошибка при получении поста",

			usecase: get.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					id int,
				) (*entity.Post, error) {
					assert.Equal(t, 1, id)

					return nil, assert.AnError
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/post/1",
					nil,
				)
				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("post_id", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				assert.JSONEq(
					t,
					`{"errors": [{"status":1003, "title": "ошибка при получении поста", "detail": "assert.AnError general error for testing"}]}`,
					resp.Body.String(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				resp := httptest.NewRecorder()

				NewHandler(
					tt.usecase,
					opentracing.NoopTracer{},
				).Handle(resp, tt.request())

				tt.want(t, resp)
			},
		)
	}
}
