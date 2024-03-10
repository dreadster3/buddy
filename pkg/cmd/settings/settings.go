package settings

import (
	"log/slog"

	"github.com/dreadster3/buddy/pkg/config"
	"github.com/dreadster3/buddy/pkg/log"
)

type Settings struct {
	Version      string
	GlobalConfig *config.GlobalConfig
	Logger       *slog.Logger
}

func New(version string, globalConfig *config.GlobalConfig) *Settings {
	return &Settings{
		Version:      version,
		GlobalConfig: globalConfig,
		Logger:       log.Logger,
	}
}
