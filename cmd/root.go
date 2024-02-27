package cmd

import (
	"github.com/spf13/cobra"
)

var buddyCmd = &cobra.Command{
	Use:   "buddy",
	Short: "buddy is a CLI tool to help you automate your development workflow",
	// Get version from buddy.json
	Version: "0.0.1-beta02",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	if err := buddyCmd.Execute(); err != nil {
		return err
	}

	return nil
}
