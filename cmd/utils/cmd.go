package utils

import (
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/cmd/utils/migrate"
)

var Cmd = &cobra.Command{
	Use:   "utils",
	Short: "вспомогательные команды",
	Long:  "вспомогательные команды для работы с базой данных, rabbitmq и тд",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(migrate.Cmd)
}
