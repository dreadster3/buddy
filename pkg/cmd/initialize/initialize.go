package initialize

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/dreadster3/buddy/pkg/template"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	Settings *settings.Settings

	// Args
	ProjectName  string
	Description  string
	TemplateName string
}

func NewCmdInit(settings *settings.Settings) *cobra.Command {
	opts := &InitOptions{
		Settings: settings,
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
			if opts.TemplateName != "" {
				templatesPath := path.Join(opts.Settings.GlobalConfig.GetTemplatesPath(), opts.TemplateName)
				if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
					opts.Settings.Logger.Error("Template does not exist", "template", opts.TemplateName)
					return errors.New("Template does not exist")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Settings.Logger.Debug("Directory provided", "directory", args[0])

				directoryPath := path.Join(opts.Settings.WorkingDirectory, args[0])
				if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
					opts.Settings.Logger.Debug("Creating directory", "directory", directoryPath)
					err := os.MkdirAll(directoryPath, 0755)
					if err != nil {
						return err
					}
				}

				opts.Settings.WorkingDirectory = directoryPath
			}

			realPath, err := filepath.Abs(opts.Settings.WorkingDirectory)
			if err != nil {
				return err
			}

			opts.ProjectName = path.Base(realPath)

			if _, err := os.Stat(path.Join(opts.Settings.WorkingDirectory, opts.Settings.GlobalConfig.FileName)); err == nil {
				opts.Settings.Logger.Error("File already exists", "file", opts.Settings.GlobalConfig.FileName)
				return errors.New("File already exists")
			}

			opts.Settings.Logger = opts.Settings.Logger.With("projectName", opts.ProjectName)

			return RunInit(opts)
		},
	}

	initCmd.Flags().StringVar(&opts.Description, "description", "A new buddy project", "Description of the project")
	initCmd.Flags().StringVarP(&opts.TemplateName, "template", "t", "", "Template to use for the project")

	initCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		entries, err := os.ReadDir(opts.Settings.GlobalConfig.GetTemplatesPath())
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		templates := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				templates = append(templates, entry.Name())
			}
		}

		return templates, cobra.ShellCompDirectiveNoFileComp
	})

	return initCmd
}

func RunInit(opts *InitOptions) error {
	opts.Settings.ProjectConfig = config.NewProjectConfig(opts.ProjectName, "0.0.1", opts.Description, opts.Settings.GlobalConfig.Author, map[string]string{})

	opts.Settings.Logger.Debug("Creating buddy file", "file", opts.Settings.GlobalConfig.FileName)
	err := opts.Settings.ProjectConfig.WriteToFile(path.Join(opts.Settings.WorkingDirectory, opts.Settings.GlobalConfig.FileName))
	if err != nil {
		return err
	}

	if opts.TemplateName != "" {
		opts.Settings.Logger.Debug("Initializing template", "template", opts.TemplateName)
		fmt.Fprintln(opts.Settings.StdOut, "Initializing template", opts.TemplateName)

		templatePath := path.Join(opts.Settings.GlobalConfig.GetTemplatesPath(), opts.TemplateName)
		projectData := template.NewProjectTemplateData(opts.Settings)

		opts.Settings.Logger.Debug("Rendering project template", "template", opts.TemplateName)
		err := template.RenderProject(opts.Settings.WorkingDirectory, templatePath, projectData)
		if err != nil {
			return err
		}
	}

	opts.Settings.Logger.Info("Project initialized")
	fmt.Fprintln(opts.Settings.StdOut, "Project initialized successfully!")

	return nil
}
