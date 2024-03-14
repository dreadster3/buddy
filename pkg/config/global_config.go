package config

import (
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	Author        string            `json:"author"`
	FileName      string            `json:"filename"`
	TemplatesPath string            `json:"templates_path" mapstructure:"templates_path"`
	Scripts       map[string]string `json:"scripts"`

	// Keep track of where the config file was found
	// Important since can be set through environment variables
	ConfigPath string
}

func NewGlobalConfig(author string, filename string, scripts map[string]string) *GlobalConfig {
	return &GlobalConfig{
		Author:        author,
		FileName:      filename,
		Scripts:       scripts,
		TemplatesPath: "templates",
		ConfigPath:    "$HOME/.config/buddy",
	}
}

func GlobalConfigFromViper(v *viper.Viper) (*GlobalConfig, error) {
	var config GlobalConfig

	opts := func(config *mapstructure.DecoderConfig) {
	}

	err := v.Unmarshal(&config, opts)
	if err != nil {
		return nil, err
	}

	config.ConfigPath = path.Dir(v.ConfigFileUsed())

	return &config, nil
}

func (c *GlobalConfig) GetTemplatesPath() string {
	return path.Join(c.ConfigPath, c.TemplatesPath)
}
