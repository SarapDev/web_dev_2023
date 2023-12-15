package comment

import (
	"github.kolesa-team.org/backend/go-example/app/api/comment/create"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app::api::comment",
	create.Module,
)
