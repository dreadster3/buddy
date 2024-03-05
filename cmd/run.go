package cmd

import (
	"fmt"

	"github.com/dreadster3/buddy/models"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	// Args
	CommandName string
	CommandArgs []string

	// Flags
	ListCommands bool
}

func init() {
	buddyCmd.AddCommand(NewCmdRun())
}

func NewCmdRun() *cobra.Command {
	opts := &RunOptions{}

	var runCmd = &cobra.Command{
		Use:                   "run [options] [command] [args...]",
		DisableFlagsInUseLine: true,
		Short:                 "Run a predefined command",
		Long:                  `Run a predefined command`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || opts.ListCommands {
				return nil
			}

			buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
			if err != nil {
				return err
			}

			if _, ok := buddyConfig.Scripts[args[0]]; !ok {
				return fmt.Errorf("Command %s not found", args[0])
			}

			return nil
		},
		Aliases: []string{"execute"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			var commands []string

			for commandName := range buddyConfig.Scripts {
				commands = append(commands, commandName)
			}

			return commands, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				opts.ListCommands = true
			} else {
				opts.CommandName = args[0]
				opts.CommandArgs = args[1:]
			}

			return runExecute(opts)
		},
	}

	runCmd.Flags().BoolVarP(&opts.ListCommands, "list", "l", false, "List all available commands")

	return runCmd
}

func runExecute(opts *RunOptions) error {
	buddyConfig, err := models.ParseBuddyConfigFile("buddy.json")
	if err != nil {
		return err
	}

	if opts.ListCommands {
		for commandName := range buddyConfig.Scripts {
			fmt.Println(commandName, "->", buddyConfig.Scripts[commandName])
		}

		return nil
	}

	return buddyConfig.RunScriptArgs(opts.CommandName, opts.CommandArgs)
}
