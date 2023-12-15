package external

import (
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/app/common/repository/post"
	"github.kolesa-team.org/backend/go-example/app/common/usecase/post/create"
	"github.kolesa-team.org/backend/go-example/app/external"
	"github.kolesa-team.org/backend/go-module/chi"
	"github.kolesa-team.org/backend/go-module/config"
	"github.kolesa-team.org/backend/go-module/env"
	"github.kolesa-team.org/backend/go-module/logger"
	"github.kolesa-team.org/backend/go-module/mysql"
	"github.kolesa-team.org/backend/go-module/rabbitmq/connection"
	"github.kolesa-team.org/backend/go-module/rabbitmq/consumer"
	"github.kolesa-team.org/backend/go-module/tracer"
	"go.uber.org/fx"
)

var (
	configFile *string
	Cmd        = &cobra.Command{
		Use:   "external",
		Short: "запуск воркера создания постов из очереди",
		Long:  "запуск воркера external для создания постов из очереди rabbitmq",
		Run: func(cmd *cobra.Command, args []string) {
			fx.New(
				// конфигурация приложения
				config.ProvideConfigPath(*configFile),
				config.ModuleWithAppName("example::worker::external"),

				// базовые модули
				logger.Module,
				env.Module,
				tracer.Module,
				chi.Module,
				mysql.Module,

				// usecase и repo модули для работы с постами
				create.Module,
				post.Module,

				// модуль соединения с rabbitmq
				connection.Module,
				// модуль консьюмера
				consumer.Module,

				// модуль приложения воркера
				external.Module,
			).Run()
		},
	}
)

func init() {
	configFile = Cmd.Flags().StringP(
		"config", "c",
		"config/development/external.toml",
		"путь к файлу конфигурации",
	)
}
