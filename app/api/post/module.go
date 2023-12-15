package post

import (
	"github.kolesa-team.org/backend/go-example/app/api/post/comments"
	"github.kolesa-team.org/backend/go-example/app/api/post/create"
	"github.kolesa-team.org/backend/go-example/app/api/post/get"
	"github.kolesa-team.org/backend/go-example/app/api/post/list"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app::api::post",
	create.Module,
	get.Module,
	list.Module,
	comments.Module,
)
