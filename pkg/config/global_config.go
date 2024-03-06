package config

import (
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	Author   string            `json:"author"`
	FileName string            `json:"filename"`
	Scripts  map[string]string `json:"scripts"`
}

func NewGlobalConfig(author string, filename string, scripts map[string]string) *GlobalConfig {
	return &GlobalConfig{
		Author:   author,
		FileName: filename,
		Scripts:  scripts,
	}
}

func GlobalConfigFromViper(v *viper.Viper) (*GlobalConfig, error) {
	var config GlobalConfig

	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
