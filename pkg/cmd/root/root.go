package root

import (
	"os"
	"path/filepath"

	"github.com/dreadster3/buddy/pkg/cmd/get"
	"github.com/dreadster3/buddy/pkg/cmd/initialize"
	"github.com/dreadster3/buddy/pkg/cmd/run"
	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/dreadster3/buddy/pkg/utils"
	"github.com/spf13/cobra"
)

type RootOptions struct {
	Settings *settings.Settings

	// Args
	CommandName string
	CommandArgs []string
}

func NewRootCmd(settings *settings.Settings) *cobra.Command {
	opts := &RootOptions{
		Settings: settings,

		CommandArgs: []string{},
	}

	var rootCmd = &cobra.Command{
		Use:                   "buddy [options] [command]",
		DisableFlagsInUseLine: true,
		Short:                 "buddy is a CLI tool to help you automate your development workflow",
		Version:               settings.Version,
		Args:                  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			subcommands := []string{}
			commandNames := []string{}

			for _, cmd := range cmd.Root().Commands() {
				subcommands = append(subcommands, cmd.Name())
			}

			if opts.Settings.ProjectConfig == nil {
				return subcommands, cobra.ShellCompDirectiveNoFileComp
			}

			for commandName := range opts.Settings.ProjectConfig.Scripts {
				commandNames = append(commandNames, commandName)
			}

			return utils.SetDifference(commandNames, subcommands), cobra.ShellCompDirectiveNoFileComp
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			absoluteWorkingDir, err := filepath.Abs(opts.Settings.WorkingDirectory)
			if err != nil {
				return err
			}

			opts.Settings.WorkingDirectory = absoluteWorkingDir
			opts.Settings.Logger = opts.Settings.Logger.With("workingDirectory", opts.Settings.WorkingDirectory)

			projectConfig, err := config.ParseMergeProjectConfigFile(opts.Settings.GlobalConfig)
			if err != nil {
				projectConfig = config.NewProjectConfig("buddy", opts.Settings.Version, "", opts.Settings.GlobalConfig.Author, opts.Settings.GlobalConfig.Scripts)
			}

			opts.Settings.ProjectConfig = projectConfig

			return os.Chdir(opts.Settings.WorkingDirectory)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CommandName = args[0]

			if len(args) > 1 {
				opts.CommandArgs = args[1:]
			}

			opts.Settings.Logger = opts.Settings.Logger.With("command", opts.CommandName, "args", opts.CommandArgs)

			return RunRoot(opts)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&opts.Settings.WorkingDirectory, "directory", "d", ".", "The working directory to run the command in")

	rootCmd.AddCommand(initialize.NewCmdInit(opts.Settings))
	rootCmd.AddCommand(get.NewCmdGet(opts.Settings))
	rootCmd.AddCommand(run.NewCmdRun(opts.Settings))

	return rootCmd
}

func RunRoot(opts *RootOptions) error {
	opts.Settings.Logger.Debug("No built in command found, running custom command")

	runOpts := &run.RunOptions{
		Settings: opts.Settings,

		CommandName: opts.CommandName,
		CommandArgs: opts.CommandArgs,
	}

	return run.RunExecute(runOpts)
}
