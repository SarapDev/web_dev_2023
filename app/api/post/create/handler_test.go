package create

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/create"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name    string
		usecase create.PostCreateExecutor
		request func() *http.Request
		want    func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "успешное создание поста",
			usecase: create.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					title string,
					content string,
				) (*entity.Post, error) {
					return &entity.Post{
						ID:      1,
						Title:   title,
						Content: content,
					}, nil
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodPost,
					"/post",
					strings.NewReader(
						`{"title":"hello world","content":"hi my friend. how are you?"}`,
					),
				)

				req.Header.Set("Content-Type", "application/json")

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, resp.Code)
				assert.JSONEq(
					t,
					`{"data": {"id":1,"title":"hello world","content":"hi my friend. how are you?"}}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "ошибка при создании поста",
			usecase: create.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					title string,
					content string,
				) (*entity.Post, error) {
					return nil, errors.New("some error")
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodPost,
					"/post",
					strings.NewReader(
						`{"title":"hello world","content":"hi my friend. how are you?"}`,
					),
				)

				req.Header.Set("Content-Type", "application/json")

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				assert.JSONEq(
					t,
					`
						{
						  "errors": [
							{
							  "status": 1003,
							  "title": "internal server error",
								"detail": "some error"
							}
						  ]
						}
					`,
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
