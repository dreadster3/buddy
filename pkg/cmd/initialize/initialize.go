package initialize

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	Settings *settings.Settings

	// Args
	ProjectName string

	// Flags
	WorkingDirectory string
}

func NewCmdInit(settings *settings.Settings) *cobra.Command {
	opts := &InitOptions{
		Settings: settings,

		WorkingDirectory: ".",
	}

	var initCmd = &cobra.Command{
		Use:                   "init [options] [directory]",
		Version:               settings.Version,
		DisableFlagsInUseLine: true,
		Short:                 "Initialize a new buddy file",
		Long: `Initialize a new buddy file.
	If a directory is provided, buddy.json will be created in the directory.
	If no directory is provided, buddy.json will be created in the current directory.`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			workingDirectory, err := cmd.Flags().GetString("directory")
			if err != nil {
				return err
			}

			opts.WorkingDirectory = workingDirectory
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.WorkingDirectory = path.Join(opts.WorkingDirectory, args[0])
			}

			realPath, err := filepath.Abs(opts.WorkingDirectory)
			if err != nil {
				return err
			}

			opts.ProjectName = path.Base(realPath)

			projectConfigFilePath := path.Join(opts.WorkingDirectory, opts.Settings.GlobalConfig.FileName)
			if _, err := os.Stat(projectConfigFilePath); err == nil {
				opts.Settings.Logger.Error("File already exists", "path", projectConfigFilePath)
				return fmt.Errorf("%s already exists", opts.Settings.GlobalConfig.FileName)
			}

			opts.Settings.Logger = opts.Settings.Logger.With("project", opts.ProjectName, "projectName", opts.ProjectName)

			return RunInit(opts)
		},
	}

	return initCmd
}

func RunInit(opts *InitOptions) error {
	projectConfig := config.NewProjectConfig(opts.ProjectName, "0.0.1", "A new buddy project", opts.Settings.GlobalConfig.Author, map[string]string{})

	err := projectConfig.WriteToFile(path.Join(opts.WorkingDirectory, opts.Settings.GlobalConfig.FileName))
	if err != nil {
		return err
	}

	opts.Settings.Logger.Info("Project initialized")
	fmt.Fprintln(opts.Settings.StdOut, opts.Settings.GlobalConfig.FileName, "created")

	return nil
}
