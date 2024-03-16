package template

import (
	"time"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
)

type ProjectTemplateData struct {
	ProjectConfig *config.ProjectConfig

	DirectoryName string
	Time          time.Time
}

func NewProjectTemplateData(settings *settings.Settings) *ProjectTemplateData {
	return &ProjectTemplateData{
		ProjectConfig: settings.ProjectConfig,

		DirectoryName: settings.WorkingDirectory,
		Time:          time.Now(),
	}
}
