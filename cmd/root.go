package cmd

import (
	"github.com/dreadster3/buddy/models"
	"github.com/spf13/cobra"
)

var buddyCmd = &cobra.Command{
	Use:                   "buddy [options] [command]",
	DisableFlagsInUseLine: true,
	Short:                 "buddy is a CLI tool to help you automate your development workflow",
	Version:               "0.0.1-beta03",
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commandName := args[0]

		buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
		if err != nil {
			return err
		}

		return buddyConfig.RunScript(commandName)
	},
}

func Execute() error {
	if err := buddyCmd.Execute(); err != nil {
		return err
	}

	return nil
}
