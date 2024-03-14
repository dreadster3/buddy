package initialize

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/dreadster3/buddy/pkg/log"
)

type ProjectData struct {
	ProjectConfig *config.ProjectConfig

	Time time.Time
}

type ProjectTemplate struct {
	TemplateDirectory string
	Name              string
}

func NewProjectTemplate(templateDirectory string, name string) *ProjectTemplate {
	return &ProjectTemplate{
		TemplateDirectory: templateDirectory,
		Name:              name,
	}
}

func (t *ProjectTemplate) RenderProject(settings *settings.Settings) error {
	err := filepath.WalkDir(t.TemplateDirectory, func(objectPath string, d fs.DirEntry, err error) error {
		relativePath, err := filepath.Rel(t.TemplateDirectory, objectPath)
		if err != nil {
			return err
		}

		if relativePath == "." {
			return nil
		}

		log.Logger.Info("Walking", "path", objectPath, "dir", d, "err", err)

		if d.IsDir() {
			err := os.MkdirAll(path.Join(settings.WorkingDirectory, relativePath), 0755)
			if err != nil && !os.IsExist(err) {
				return err
			}

			return nil
		}

		filePath := path.Join(settings.WorkingDirectory, relativePath)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		return t.RenderFile(file, objectPath, settings)
	})

	return err
}

func (t *ProjectTemplate) RenderFile(writer io.Writer, filePath string, settings *settings.Settings) error {
	tmpl, err := template.New(path.Base(filePath)).ParseFiles(filePath)
	if err != nil {
		return err
	}

	projectData := &ProjectData{
		ProjectConfig: settings.ProjectConfig,

		Time: time.Now(),
	}

	err = tmpl.Execute(writer, projectData)

	return nil
}
