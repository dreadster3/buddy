package initialize

import (
	"fmt"
	"os"
	"path"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	GlobalConfig *config.GlobalConfig

	Directory   string
	ProjectName string
}

func NewCmdInit(globalConfig *config.GlobalConfig) *cobra.Command {
	opts := &InitOptions{
		Directory:    ".",
		GlobalConfig: globalConfig,
	}

	var initCmd = &cobra.Command{
		Use:                   "init [options] [directory]",
		DisableFlagsInUseLine: true,
		Short:                 "Initialize a new buddy file",
		Long: `Initialize a new buddy file.
	If a directory is provided, buddy.json will be created in the directory.
	If no directory is provided, buddy.json will be created in the current directory.`,
		Args: cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				directoryName := args[0]

				err := os.Mkdir(directoryName, 0755)
				if err != nil && !os.IsExist(err) {
					return err
				}

				opts.Directory = directoryName
				opts.ProjectName = path.Base(directoryName)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			buddyFilePath := path.Join(opts.Directory, opts.GlobalConfig.FileName)
			if _, err := os.Stat(buddyFilePath); err == nil {
				return fmt.Errorf("%s already exists", opts.GlobalConfig.FileName)
			}

			if opts.ProjectName == "" {
				workingDir, err := os.Getwd()
				if err != nil {
					return err
				}

				opts.Directory = workingDir
				opts.ProjectName = path.Base(workingDir)
			}

			return runInit(opts)
		},
	}

	return initCmd
}

func runInit(opts *InitOptions) error {
	projectConfig := config.NewProjectConfig(opts.ProjectName, "0.0.1", "A new buddy project", opts.GlobalConfig.Author, map[string]string{})

	err := projectConfig.WriteToFile(path.Join(opts.Directory, opts.GlobalConfig.FileName))
	if err != nil {
		return err
	}

	fmt.Println(opts.GlobalConfig.FileName, "created")

	return nil
}
