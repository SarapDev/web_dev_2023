package api

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.kolesa-team.org/backend/go-example/app/api"
	"github.kolesa-team.org/backend/go-example/app/common/repository/comment"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/config"
	"github.kolesa-team.org/backend/go-module/env"
	"github.kolesa-team.org/backend/go-module/logger"
	"github.kolesa-team.org/backend/go-module/mysql"
	"github.kolesa-team.org/backend/go-module/tracer"

	postUsecaseCreate "github.kolesa-team.org/backend/go-example/app/common/usecase/post/create"
	postUsecaseGet "github.kolesa-team.org/backend/go-example/app/common/usecase/post/get"
	postUsecaseList "github.kolesa-team.org/backend/go-example/app/common/usecase/post/list"

	commentUsecaseCreate "github.kolesa-team.org/backend/go-example/app/common/usecase/comment/create"
	commentUsecaseList "github.kolesa-team.org/backend/go-example/app/common/usecase/comment/list"
)

var (
	configFile *string
	Cmd        = &cobra.Command{
		Use:     "api",
		Short:   "запуск сервера API",
		Long:    "запуск HTTP/GRPC сервера Rodina API",
		Example: "binary serve api --config config/development/api.toml",
		Run: func(cmd *cobra.Command, _ []string) {
			fx.New(
				// конфигурация приложения
				config.ProvideConfigPath(*configFile),
				config.ModuleWithAppName("example::api"),

				// подключение нужных модулей
				logger.Module,
				env.Module,
				tracer.Module,
				chi.Module,
				mysql.Module,

				// регистрация репозиториев
				post.Module,
				comment.Module,

				// регистрация usecase -ов
				postUsecaseCreate.Module,
				postUsecaseGet.Module,
				postUsecaseList.Module,
				commentUsecaseCreate.Module,
				commentUsecaseList.Module,

				// подключение модуля приложения API
				api.Module,
			).Run()
		},
	}
)

func init() {
	configFile = Cmd.Flags().StringP(
		"config", "c",
		"config/development/api.toml",
		"путь к конфигурационному файлу",
	)
}
