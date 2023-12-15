package create

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/comment/create"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name    string
		usecase create.CommentCreateExecutor
		request func() *http.Request
		want    func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "успешное создание комментария",
			usecase: create.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					postID int,
					author string,
					content string,
				) (*entity.Comment, error) {
					assert.Equal(t, 1, postID)
					assert.Equal(t, "author", author)
					assert.Equal(t, "hi my friend. how are you?", content)

					return &entity.Comment{
						ID:      1,
						PostID:  postID,
						Author:  author,
						Content: content,
					}, nil
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodPost,
					"/comment",
					strings.NewReader(
						`{"post_id":1,"author":"author","content":"hi my friend. how are you?"}`,
					),
				)

				req.Header.Set("Content-Type", "application/json")

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, resp.Code)
				assert.JSONEq(
					t,
					`{"data": {"id":1,"post_id":1,"author":"author","content":"hi my friend. how are you?"}}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "ошибка валидации",
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodPost,
					"/comment",
					strings.NewReader(
						`{"post_id":1,"author":"author","content":""}`,
					),
				)

				req.Header.Set("Content-Type", "application/json")

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				assert.JSONEq(
					t,
					`{"errors": [{"status":1002, "title":"validation error", "detail":"content: cannot be blank"}]}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "невалидный JSON",
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodPost,
					"/comment",
					strings.NewReader(
						`{"post_id":1,"author":"author","content":"hi my friend. how are you?"`,
					),
				)

				req.Header.Set("Content-Type", "application/json")

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
				assert.JSONEq(
					t,
					`{"errors": [{"status":1001, "title":"ошибка при чтении тела запроса", "detail":"unexpected end of JSON input"}]}`,
					resp.Body.String(),
				)
			},
		},
		{
			name: "ошибка во время выполнения usecase",

			usecase: create.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					postID int,
					author string,
					content string,
				) (*entity.Comment, error) {
					return nil, assert.AnError
				},
			},
			request: func() *http.Request {
				req := httptest.NewRequest(
					http.MethodPost,
					"/comment",
					strings.NewReader(
						`{"post_id":1,"author":"author","content":"hi my friend. how are you?"}`,
					),
				)

				req.Header.Set("Content-Type", "application/json")

				return req
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				assert.JSONEq(
					t,
					`{"errors": [{"status":1004, "title":"ошибка при создании комментария", "detail":"assert.AnError general error for testing"}]}`,
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
