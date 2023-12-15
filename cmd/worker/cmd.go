package worker

import (
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/cmd/worker/external"
)

var Cmd = &cobra.Command{
	Use:   "worker",
	Short: "запуск воркера",
	Long:  "запуск доступного воркера для обработки задач из очереди",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(external.Cmd)
}
