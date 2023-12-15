package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.kolesa-team.org/backend/go-example/cmd/serve"
	"github.kolesa-team.org/backend/go-example/cmd/utils"
	"github.kolesa-team.org/backend/go-example/cmd/worker"
)

var (
	rootCmd = &cobra.Command{
		Version: "0.1.0",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Help()
		},
	}
)

func init() {
	// скрываем дефолтную команду `completion`
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// подключаем дочерние команды
	rootCmd.AddCommand(serve.Cmd)
	rootCmd.AddCommand(worker.Cmd)
	rootCmd.AddCommand(utils.Cmd)
}

// Execute запускает корневую команду
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
