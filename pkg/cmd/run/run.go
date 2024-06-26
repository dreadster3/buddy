package run

import (
	"fmt"
	"text/tabwriter"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/spf13/cobra"
)

type RunOptions struct {
	// Global Config
	Settings *settings.Settings

	// Args
	ScriptName string
	ScriptArgs []string

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
		Args:                  cobra.ArbitraryArgs,
		Aliases:               []string{"execute"},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			projectConfig := opts.Settings.ProjectConfig

			var scripts []string

			for scriptName := range projectConfig.Scripts {
				scripts = append(scripts, fmt.Sprintf("%s\t%s", scriptName, projectConfig.Scripts[scriptName]))
			}

			return scripts, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				opts.ListCommands = true
				opts.Settings.Logger.Info("No command provided, listing all commands")
			} else {
				opts.ScriptName = args[0]
				opts.ScriptArgs = args[1:]
				opts.Settings.Logger = opts.Settings.Logger.With("command", opts.ScriptName, "args", opts.ScriptArgs)
			}

			return RunExecute(opts)
		},
	}

	runCmd.Flags().BoolVarP(&opts.ListCommands, "list", "l", false, "List all available commands")

	return runCmd
}

func RunExecute(opts *RunOptions) error {
	opts.Settings.Logger.Debug("Executing Script")

	if opts.ListCommands {
		writer := tabwriter.NewWriter(opts.Settings.StdOut, 0, 0, 1, ' ', 0)
		for commandName := range opts.Settings.ProjectConfig.Scripts {
			fmt.Fprintln(writer, commandName, "\t->\t", opts.Settings.ProjectConfig.Scripts[commandName])
		}
		writer.Flush()

		return nil
	}

	err := opts.Settings.ProjectConfig.RunScriptArgs(opts.ScriptName, opts.ScriptArgs)
	if err != nil {
		return err
	}

	return nil
}
