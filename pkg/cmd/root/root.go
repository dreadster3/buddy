package root

import (
	"os"

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
	ScriptName string
	ScriptArgs []string
}

func NewRootCmd(settings *settings.Settings) *cobra.Command {
	opts := &RootOptions{
		Settings: settings,

		ScriptArgs: []string{},
	}

	var rootCmd = &cobra.Command{
		Use:                   "buddy [options] [command]",
		DisableFlagsInUseLine: true,
		Short:                 "buddy is a CLI tool to help you automate your development workflow",
		Version:               settings.Version,
		Args:                  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			subCommands := []string{}
			scriptNames := []string{}

			for _, cmd := range cmd.Root().Commands() {
				subCommands = append(subCommands, cmd.Name())
			}

			if opts.Settings.ProjectConfig == nil {
				return subCommands, cobra.ShellCompDirectiveNoFileComp
			}

			for scriptName := range opts.Settings.ProjectConfig.Scripts {
				scriptNames = append(scriptNames, scriptName)
			}

			return utils.SetDifference(scriptNames, subCommands), cobra.ShellCompDirectiveNoFileComp
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Settings.Logger = opts.Settings.Logger.With("workingDirectory", opts.Settings.WorkingDirectory)

			projectConfig, err := config.ParseMergeProjectConfigFile(opts.Settings.GlobalConfig)
			if err != nil {
				projectConfig = config.NewProjectConfig("buddy", opts.Settings.Version, "", opts.Settings.GlobalConfig.Author, opts.Settings.GlobalConfig.Scripts)
			}

			opts.Settings.ProjectConfig = projectConfig

			err = os.Chdir(opts.Settings.WorkingDirectory)
			if err != nil {
				return err
			}

			opts.Settings.WorkingDirectory = "."
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ScriptName = args[0]

			if len(args) > 1 {
				opts.ScriptArgs = args[1:]
			}

			opts.Settings.Logger = opts.Settings.Logger.With("command", opts.ScriptName, "args", opts.ScriptArgs)

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

		ScriptName: opts.ScriptName,
		ScriptArgs: opts.ScriptArgs,
	}

	return run.RunExecute(runOpts)
}
