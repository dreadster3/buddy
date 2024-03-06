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
}

func NewRootCmd(globalConfig *config.GlobalConfig) *cobra.Command {
	opts := &RootOptions{
		GlobalConfig: globalConfig,
	}

	var buddyCmd = &cobra.Command{
		Use:                   "buddy [options] [command]",
		DisableFlagsInUseLine: true,
		Short:                 "buddy is a CLI tool to help you automate your development workflow",
		Version:               "v0.0.1-beta.4",
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.CommandName = args[0]

			return runRoot(opts)
		},
	}

	buddyCmd.AddCommand(initialize.NewCmdInit(opts.GlobalConfig))
	buddyCmd.AddCommand(get.NewCmdGet(opts.GlobalConfig))
	buddyCmd.AddCommand(run.NewCmdRun(opts.GlobalConfig))

	return buddyCmd
}

func runRoot(opts *RootOptions) error {
	projectConfig, err := config.ParseMergeProjectConfigFile(opts.GlobalConfig)
	if err != nil {
		return err
	}

	return projectConfig.RunScript(opts.CommandName)
}
