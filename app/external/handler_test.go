package external

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/wagslane/go-rabbitmq"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/create"
	"golang.org/x/exp/slog"
	"io"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name     string
		usecase  create.PostCreateExecutor
		delivery func() rabbitmq.Delivery
		want     func(t *testing.T, action rabbitmq.Action)
	}{
		{
			name: "успешное создание поста",
			usecase: create.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					title, content string,
				) (*entity.Post, error) {
					return &entity.Post{
						ID:      1,
						Title:   title,
						Content: content,
					}, nil
				},
			},
			delivery: func() rabbitmq.Delivery {
				return rabbitmq.Delivery{
					Delivery: amqp091.Delivery{
						Body: []byte(`{"title":"hello world","content":"hi my friend. how are you?"}`),
					},
				}
			},
			want: func(t *testing.T, action rabbitmq.Action) {
				assert.Equal(t, rabbitmq.Ack, action)
			},
		},
		{
			name: "ошибка создания поста",
			usecase: create.MockExecutor{
				MockExecuteFunc: func(
					ctx context.Context,
					title, content string,
				) (*entity.Post, error) {
					return nil, assert.AnError
				},
			},
			delivery: func() rabbitmq.Delivery {
				return rabbitmq.Delivery{
					Delivery: amqp091.Delivery{
						Body: []byte(`{"title":"hello world","content":"hi my friend. how are you?"}`),
					},
				}
			},
			want: func(t *testing.T, action rabbitmq.Action) {
				assert.Equal(t, rabbitmq.NackRequeue, action)
			},
		},
		{
			name: "ошибка парсинга",
			delivery: func() rabbitmq.Delivery {
				return rabbitmq.Delivery{
					Delivery: amqp091.Delivery{
						Body: []byte(`{"title":"hello world","content":"hi my friend. how are you?}`),
					},
				}
			},
			want: func(t *testing.T, action rabbitmq.Action) {
				assert.Equal(t, rabbitmq.NackDiscard, action)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.want(
					t, NewHandler(
						slog.New(slog.NewJSONHandler(io.Discard, nil)),
						tt.usecase,
						opentracing.NoopTracer{},
					).Handle(
						opentracing.ContextWithSpan(
							context.Background(),
							opentracing.StartSpan("test"),
						),
						tt.delivery(),
					),
				)
			},
		)
	}
}
