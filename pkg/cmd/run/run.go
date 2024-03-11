package run

import (
	"fmt"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	// Global Config
	Settings *settings.Settings

	// Args
	CommandName string
	CommandArgs []string

	// Flags
	ListCommands bool
}

func NewCmdRun(settings *settings.Settings) *cobra.Command {
	opts := &RunOptions{
		Settings: settings,
	}

	var runCmd = &cobra.Command{
		Use:                   "run [options] [command] [args...]",
		DisableFlagsInUseLine: true,
		Version:               settings.Version,
		Short:                 "Run a predefined command",
		Long:                  `Run a predefined command`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || opts.ListCommands {
				return nil
			}

			if _, ok := opts.Settings.ProjectConfig.Scripts[args[0]]; !ok {
				return fmt.Errorf("Command %s not found", args[0])
			}

			return nil
		},
		Aliases: []string{"execute"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			projectConfig := opts.Settings.ProjectConfig

			var commands []string

			for commandName := range projectConfig.Scripts {
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
				opts.Settings.Logger = opts.Settings.Logger.With("command", opts.CommandName, "args", opts.CommandArgs)
			}

			return RunExecute(opts)
		},
	}

	runCmd.Flags().BoolVarP(&opts.ListCommands, "list", "l", false, "List all available commands")

	return runCmd
}

func RunExecute(opts *RunOptions) error {
	opts.Settings.Logger.Debug("Executing Command")

	if opts.ListCommands {
		opts.Settings.Logger.Info("No command provided, listing all commands")
		for commandName := range opts.Settings.ProjectConfig.Scripts {
			fmt.Println(commandName, "->", opts.Settings.ProjectConfig.Scripts[commandName])
		}

		return nil
	}

	return opts.Settings.ProjectConfig.RunScriptArgs(opts.CommandName, opts.CommandArgs)
}
