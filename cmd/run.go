package cmd

import (
	"bytes"
	"fmt"
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
	Use:   "run",
	Short: "Run a predefined command",
	Long:  `Run a predefined command`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if listCommands || len(args) == 0 {
			buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			for commandName := range buddyConfig.Scripts {
				cmd.Println(commandName, "->", buddyConfig.Scripts[commandName])
			}

			return
		}

		commandName := args[0]

		buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		command, ok := buddyConfig.Scripts[commandName]
		if !ok {
			cmd.PrintErrf("Command %s not found", commandName)
			return
		}

		var stdOutBuf, stdErrBuf bytes.Buffer
		execCommand := exec.Command("sh", "-c", command)
		execCommand.Stdout = &stdOutBuf
		execCommand.Stderr = &stdErrBuf

		err = execCommand.Run()
		if err != nil {
			cmd.PrintErrf("Error: %s", stdErrBuf.String())
			return
		}

		output := stdOutBuf.Bytes()

		if len(output) > 0 {
			fmt.Println(string(output))
		}
	},
}
