package migrate

import (
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/cmd/utils/migrate/db"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "миграции",
	Long:  "команда для миграции базы данных и не только",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(
		db.Cmd,
	)
}
