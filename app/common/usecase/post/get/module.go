package get

import (
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

var Module = fx.Module(
	"app::common::usecase::post::get",
	fx.Decorate(
		func(logger *slog.Logger) *slog.Logger {
			return logger.With(
				"module",
				"app::common::usecase::post::get",
			)
		},
	),
	fx.Provide(
		fx.Annotate(
			NewExecutor,
			fx.As(new(PostGetExecutor)),
		),
	),
)
