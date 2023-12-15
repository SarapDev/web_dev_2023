package db

import (
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/cmd/utils/migrate/db/down"
	"github.kolesa-team.org/backend/go-example/cmd/utils/migrate/db/up"
)

var Cmd = &cobra.Command{
	Use:   "db",
	Short: "миграции для базы данных",
	Long:  "команда для миграции базы данных",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.PersistentFlags().StringP(
		"url", "u",
		"",
		"url для подключения к mysql",
	)

	Cmd.AddCommand(up.Cmd)
	Cmd.AddCommand(down.Cmd)
}
