package initialize

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	directory   string
	projectName string
}

func NewCmdInit() *cobra.Command {
	opts := &InitOptions{
		directory: ".",
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

				opts.directory = directoryName
				opts.projectName = path.Base(directoryName)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			buddyFilePath := path.Join(opts.directory, "buddy.json")
			if _, err := os.Stat(buddyFilePath); err == nil {
				return errors.New("buddy.json already exists")
			}

			if opts.projectName == "" {
				workingDir, err := os.Getwd()
				if err != nil {
					return err
				}

				opts.directory = workingDir
				opts.projectName = path.Base(workingDir)
			}

			return runInit(opts)
		},
	}

	return initCmd
}

func runInit(opts *InitOptions) error {
	projectConfig := config.NewProjectConfig(opts.projectName, "0.0.1", "A new buddy project", "Anonymous", map[string]string{})

	err := projectConfig.WriteToFile(path.Join(opts.directory, "buddy.json"))
	if err != nil {
		return err
	}

	fmt.Println("buddy.json created")

	return nil
}
