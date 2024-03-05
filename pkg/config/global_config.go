package config

import "fmt"

type GlobalConfig struct {
	Author  string            `json:"author"`
	Scripts map[string]string `json:"scripts"`
}

func NewGlobalConfig(author string, scripts map[string]string) *GlobalConfig {
	return &GlobalConfig{
		Author:  author,
		Scripts: scripts,
	}
}

func ParseGlobalConfigFile(filePath string) (*GlobalConfig, error) {
	return nil, fmt.Errorf("Not implemented")
}
