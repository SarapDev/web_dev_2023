package create

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

var (
	moduleName = "app::api::post::create"
	Module     = fx.Module(
		moduleName,
		fx.Provide(NewHandler),
		fx.Invoke(
			func(router *chi.Mux, h Handler) {
				router.Post("/post", h.Handle)
			},
		),
	)
)
