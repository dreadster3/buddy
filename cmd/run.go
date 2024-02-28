package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
				cmd.Println(commandName, "->", buddyConfig.Scripts[commandName])
			}

			return nil
		}

		commandName := args[0]

		buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
		if err != nil {
			return err
		}

		command, ok := buddyConfig.Scripts[commandName]
		if !ok {
			return fmt.Errorf("Command %s not found", commandName)
		}

		execCommand := exec.Command("sh", "-c", command)
		execCommand.Stdout = os.Stdout
		execCommand.Stderr = os.Stderr

		err = execCommand.Run()
		if err != nil {
			return err
		}

		return nil
	},
}
