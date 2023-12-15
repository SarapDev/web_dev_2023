package comment

import "go.uber.org/fx"

var Module = fx.Module(
	"app::common::repository::comment",
	fx.Provide(
		fx.Annotate(
			NewMySQLRepository,
			fx.As(new(Repository)),
		),
	),
)
