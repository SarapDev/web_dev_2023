package static

import (
	"embed"
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/app/static"
	"github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/config"
	"github.kolesa-team.org/backend/go-module/env"
	"github.kolesa-team.org/backend/go-module/logger"
	"github.kolesa-team.org/backend/go-module/tracer"
	"go.uber.org/fx"
)

var (
	configFile *string
	Cmd        = &cobra.Command{
		Use:     "static",
		Short:   "запуск сервера раздачи статики",
		Long:    "запуск HTTP сервера раздачи статических файлов",
		Example: "binary serve static --config config/development/static.toml",
		Run: func(cmd *cobra.Command, _ []string) {
			fx.New(
				config.ProvideConfigPath(*configFile),
				config.ModuleWithAppName("example::static"),

				logger.Module,
				env.Module,
				tracer.Module,
				chi.Module,

				fx.Supply(
					cmd.Context().Value("fs::public").(embed.FS),
				),
				static.Module,
			).Run()
		},
	}
)

func init() {
	configFile = Cmd.Flags().StringP(
		"config", "c",
		"config/development/static.toml",
		"путь к конфигурационному файлу",
	)
}
