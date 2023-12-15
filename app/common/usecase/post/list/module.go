package list

import "go.uber.org/fx"

var Module = fx.Module(
	"app::common::usecase::post::list",
	fx.Provide(
		fx.Annotate(
			NewExecutor,
			fx.As(new(PostListExecutor)),
		),
	),
)
