package settings

import (
	"io"
	"log/slog"
	"os"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/dreadster3/buddy/pkg/log"
)

type Settings struct {
	Version          string
	GlobalConfig     *config.GlobalConfig
	ProjectConfig    *config.ProjectConfig
	Logger           *slog.Logger
	WorkingDirectory string

	StdOut io.Writer
	ErrOut io.Writer
}

func New(version string, globalConfig *config.GlobalConfig) *Settings {
	return &Settings{
		Version:          version,
		GlobalConfig:     globalConfig,
		Logger:           log.Logger,
		WorkingDirectory: ".",

		StdOut: os.Stdout,
		ErrOut: os.Stderr,
	}
}
