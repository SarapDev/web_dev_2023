package serve

import (
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/cmd/serve/api"
	"github.kolesa-team.org/backend/go-example/cmd/serve/static"
)

var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "запуск сервера",
	Long:  "запуск выбранного сервера из доступных",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(api.Cmd)
	Cmd.AddCommand(static.Cmd)
}
