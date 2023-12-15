package create

import (
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

var Module = fx.Module(
	"app::common::usecase::comment::create",
	fx.Decorate(
		func(logger *slog.Logger) *slog.Logger {
			return logger.With(
				"module",
				"app::common::usecase::comment::create",
			)
		},
	),
	fx.Provide(

		fx.Annotate(
			NewExecutor,
			fx.As(new(CommentCreateExecutor)),
		),
	),
)
