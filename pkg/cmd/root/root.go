package root

import (
	"github.com/dreadster3/buddy/pkg/cmd/get"
	"github.com/dreadster3/buddy/pkg/cmd/initialize"
	"github.com/dreadster3/buddy/pkg/cmd/run"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

type RootOptions struct {
	GlobalConfig *config.GlobalConfig

	// Args
	CommandName string
	CommandArgs []string

	// Flags
	WorkingDirectory string
}

func NewRootCmd(version string, globalConfig *config.GlobalConfig) *cobra.Command {
	opts := &RootOptions{
		GlobalConfig: globalConfig,
	}

	var rootCmd = &cobra.Command{
		Use:                   "buddy [options] [command]",
		DisableFlagsInUseLine: true,
		Short:                 "buddy is a CLI tool to help you automate your development workflow",
		Version:               version,
		Args:                  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CommandName = args[0]

			if len(args) > 1 {
				opts.CommandArgs = args[1:]
			}

			return runRoot(opts)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&opts.WorkingDirectory, "directory", "d", ".", "The working directory to run the command in")

	rootCmd.AddCommand(initialize.NewCmdInit(opts.GlobalConfig))
	rootCmd.AddCommand(get.NewCmdGet(opts.GlobalConfig))
	rootCmd.AddCommand(run.NewCmdRun(opts.GlobalConfig))

	return rootCmd
}

func runRoot(opts *RootOptions) error {
	runOpts := &run.RunOptions{
		GlobalConfig:     opts.GlobalConfig,
		CommandName:      opts.CommandName,
		CommandArgs:      opts.CommandArgs,
		WorkingDirectory: opts.WorkingDirectory,
	}

	return run.RunExecute(runOpts)
}
