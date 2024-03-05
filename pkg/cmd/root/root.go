package root

import (
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var buddyCmd = &cobra.Command{
		Use:                   "buddy [options] [command]",
		DisableFlagsInUseLine: true,
		Short:                 "buddy is a CLI tool to help you automate your development workflow",
		Version:               "dev",
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			commandName := args[0]

			buddyConfig, err := config.ParseProjectConfigFile("buddy.json")
			if err != nil {
				return err
			}

			return buddyConfig.RunScript(commandName)
		},
	}

	return buddyCmd
}
