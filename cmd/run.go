package cmd

import (
	"fmt"

	"github.com/dreadster3/buddy/models"
	"github.com/spf13/cobra"
)

var listCommands bool

func init() {
	runCmd.Flags().BoolVarP(&listCommands, "list", "l", false, "List all available commands")
	buddyCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [flags] [command]",
	Short: "Run a predefined command",
	Long:  `Run a predefined command`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if listCommands || len(args) == 0 {
			buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
			if err != nil {
				return err
			}

			for commandName := range buddyConfig.Scripts {
				fmt.Println(commandName, "->", buddyConfig.Scripts[commandName])
			}

			return nil
		}

		commandName := args[0]

		buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
		if err != nil {
			return err
		}

		return buddyConfig.RunScript(commandName)
	},
}
