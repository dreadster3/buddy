package template

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dreadster3/buddy/pkg/log"
)

func RenderTemplate(writer io.Writer, templatePath string, data any) error {
	templateName := filepath.Base(templatePath)
	tmpl, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(writer, templateName, data)
}

func RenderProject(workingDirectory string, projectTemplatePath string, data *ProjectTemplateData) error {
	err := filepath.WalkDir(projectTemplatePath, func(objectPath string, d fs.DirEntry, err error) error {
		logger := log.Logger.With("path", objectPath, "dir", d, "err", err)

		relativePath, err := filepath.Rel(projectTemplatePath, objectPath)
		if err != nil {
			return err
		}
		logger = logger.With("relativePath", relativePath)

		if relativePath == "." {
			return nil
		}

		if d.IsDir() {
			logger.Debug("Found directory, creating directory")
			err := os.MkdirAll(filepath.Join(workingDirectory, relativePath), 0755)
			if err != nil && !os.IsExist(err) {
				return err
			}

			return nil
		}

		logger.Debug("Rendering file")
		filePath := filepath.Join(workingDirectory, relativePath)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		absolutePath, err := filepath.Abs(filePath)
		if err != nil {
			return err
		}

		data.DirectoryName = filepath.Base(filepath.Dir(absolutePath))
		return RenderTemplate(file, objectPath, data)
	})

	return err
}
