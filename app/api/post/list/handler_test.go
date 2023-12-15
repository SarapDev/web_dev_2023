package list

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/list"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name    string
		usecase list.PostListExecutor
		request func() *http.Request
		want    func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "успешное получение списка постов",
			usecase: list.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					offset, limit int,
				) ([]*entity.Post, int, error) {
					return []*entity.Post{
						{
							ID:      1,
							Title:   "hello world",
							Content: "hi my friend. how are you?",
							CreatedAt: time.Date(
								2020,
								time.January,
								1,
								0,
								0,
								0,
								0,
								time.UTC,
							),
						},
					}, 10, nil
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/post",
					nil,
				)

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				assert.JSONEq(
					t,
					`{"meta":{"total":10},"data":[{"id":1,"title":"hello world","created_at":"2020-01-01T00:00:00Z"}]}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "ошибка получения списка постов",
			usecase: list.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					offset, limit int,
				) ([]*entity.Post, int, error) {
					return nil, 0, assert.AnError
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodGet,
					"/post",
					nil,
				)

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				assert.JSONEq(
					t,
					`{"errors":[{"status":1002,"title":"ошибка получения постов","detail":"assert.AnError general error for testing"}]}`,
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
