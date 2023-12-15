package post

import "go.uber.org/fx"

var Module = fx.Module(
	"app::common::repository::post",
	fx.Provide(
		fx.Annotate(
			NewMySQLRepository,
			fx.As(new(Repository)),
		),
	),
)
