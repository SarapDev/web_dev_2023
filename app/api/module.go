package api

import (
	"github.com/go-chi/chi/v5"
	"github.kolesa-team.org/backend/go-example/app/api/comment"
	"github.kolesa-team.org/backend/go-example/app/api/post"
	"github.kolesa-team.org/backend/go-module/config"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"app::api",

	// регистрируем модуль с обработчиками запросов к /post
	post.Module,
	comment.Module,

	// регистрируем обработчик запроса на /
	fx.Provide(
		fx.Annotate(
			NewRootHandler,
			fx.ParamTags(
				``,
				config.GetAppNameTag(false),
				config.GetBranchNameTag(true),
			),
		),
	),

	// регистрируем health check функцию
	// http chi модуль автоматически регистрирует обработчик запроса по пути:
	// /_private_/health
	fx.Provide(HealthHandlerFunc),

	fx.Invoke(
		func(
			lifecycle fx.Lifecycle,
			router *chi.Mux,
			rootHandler *RootHandler,
		) {
			lifecycle.Append(
				fx.StartHook(
					func() {
						router.Get("/", rootHandler.Handler)
					},
				),
			)
		},
	),
)
