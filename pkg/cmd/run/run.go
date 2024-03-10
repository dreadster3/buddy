package run

import (
	"fmt"
	"os"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	// Global Config
	Settings *settings.Settings

	// Args
	CommandName string
	CommandArgs []string

	// Flags
	ListCommands     bool
	WorkingDirectory string
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

			projectConfig, err := config.ParseMergeProjectConfigFile(opts.Settings.GlobalConfig)
			if err != nil {
				return err
			}

			if _, ok := projectConfig.Scripts[args[0]]; !ok {
				return fmt.Errorf("Command %s not found", args[0])
			}

			return nil
		},
		Aliases: []string{"execute"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			projectConfig, err := config.ParseMergeProjectConfigFile(opts.Settings.GlobalConfig)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			var commands []string

			for commandName := range projectConfig.Scripts {
				commands = append(commands, commandName)
			}

			return commands, cobra.ShellCompDirectiveNoFileComp
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			workingDirectory, err := cmd.Flags().GetString("directory")
			if err != nil {
				return err
			}

			opts.WorkingDirectory = workingDirectory
			return nil
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

	projectConfig, err := config.ParseMergeProjectConfigFile(opts.Settings.GlobalConfig)
	if err != nil {
		opts.Settings.Logger.Error("Error parsing project config file", "error", err)
		return err
	}

	if opts.ListCommands {
		opts.Settings.Logger.Info("No command provided, listing all commands")
		for commandName := range projectConfig.Scripts {
			fmt.Println(commandName, "->", projectConfig.Scripts[commandName])
		}

		return nil
	}

	err = os.Chdir(opts.WorkingDirectory)
	if err != nil {
		opts.Settings.Logger.Error("Error changing directory", "error", err)
		return err
	}

	return projectConfig.RunScriptArgs(opts.CommandName, opts.CommandArgs)
}
