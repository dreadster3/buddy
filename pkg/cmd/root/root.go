package root

import (
	"github.com/dreadster3/buddy/pkg/cmd/get"
	"github.com/dreadster3/buddy/pkg/cmd/initialize"
	"github.com/dreadster3/buddy/pkg/cmd/run"
	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/spf13/cobra"
)

type RootOptions struct {
	Settings *settings.Settings

	// Args
	CommandName string
	CommandArgs []string

	// Flags
	WorkingDirectory string
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Settings.Logger = opts.Settings.Logger.With("workingDirectory", opts.WorkingDirectory)

			return nil
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

	rootCmd.PersistentFlags().StringVarP(&opts.WorkingDirectory, "directory", "d", ".", "The working directory to run the command in")

	rootCmd.AddCommand(initialize.NewCmdInit(opts.Settings))
	rootCmd.AddCommand(get.NewCmdGet(opts.Settings))
	rootCmd.AddCommand(run.NewCmdRun(opts.Settings))

	return rootCmd
}

func RunRoot(opts *RootOptions) error {
	opts.Settings.Logger.Debug("No built in command found, running custom command")

	runOpts := &run.RunOptions{
		Settings: opts.Settings,

		CommandName:      opts.CommandName,
		CommandArgs:      opts.CommandArgs,
		WorkingDirectory: opts.WorkingDirectory,
	}

	return run.RunExecute(runOpts)
}
