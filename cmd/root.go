package cmd

import (
	"github.com/spf13/cobra"
)

var buddyCmd = &cobra.Command{
	Use:     "buddy",
	Short:   "buddy is a CLI tool to help you automate your development workflow",
	Version: "0.0.1-beta01",
}

func Execute() error {
	if err := buddyCmd.Execute(); err != nil {
		return err
	}

	return nil
}
