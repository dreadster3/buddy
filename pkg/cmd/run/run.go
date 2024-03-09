package run

import (
	"fmt"
	"os"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	// Global Config
	GlobalConfig     *config.GlobalConfig
	WorkingDirectory string

	// Args
	CommandName string
	CommandArgs []string

	// Flags
	ListCommands bool
}

func NewCmdRun(globalConfig *config.GlobalConfig) *cobra.Command {
	opts := &RunOptions{
		GlobalConfig: globalConfig,
	}

	var runCmd = &cobra.Command{
		Use:                   "run [options] [command] [args...]",
		DisableFlagsInUseLine: true,
		Short:                 "Run a predefined command",
		Long:                  `Run a predefined command`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || opts.ListCommands {
				return nil
			}

			projectConfig, err := config.ParseMergeProjectConfigFile(opts.GlobalConfig)
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
			projectConfig, err := config.ParseMergeProjectConfigFile(opts.GlobalConfig)
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
			}

			return RunExecute(opts)
		},
	}

	runCmd.Flags().BoolVarP(&opts.ListCommands, "list", "l", false, "List all available commands")

	return runCmd
}

func RunExecute(opts *RunOptions) error {
	projectConfig, err := config.ParseMergeProjectConfigFile(opts.GlobalConfig)
	if err != nil {
		return err
	}

	if opts.ListCommands {
		for commandName := range projectConfig.Scripts {
			fmt.Println(commandName, "->", projectConfig.Scripts[commandName])
		}

		return nil
	}

	err = os.Chdir(opts.WorkingDirectory)
	if err != nil {
		return err
	}

	return projectConfig.RunScriptArgs(opts.CommandName, opts.CommandArgs)
}
