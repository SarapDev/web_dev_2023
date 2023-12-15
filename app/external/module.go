package external

import (
	"github.kolesa-team.org/backend/go-module/rabbitmq/consumer"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app::external",

	fx.Provide(NewHandler),
	fx.Provide(
		func(h Handler) consumer.Handler {
			return h.Handle
		},
	),
)
