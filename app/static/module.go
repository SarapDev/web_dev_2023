package static

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"io/fs"
	"net/http"
)

var Module = fx.Module(
	"app::static",
	fx.Invoke(
		func(router *chi.Mux, eFS embed.FS) error {
			publicFS, err := fs.Sub(eFS, "public")

			if err != nil {
				return err
			}

			router.Handle("/", http.FileServer(http.FS(publicFS)))

			return nil
		},
	),
)
