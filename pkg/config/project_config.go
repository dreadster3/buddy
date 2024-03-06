package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dreadster3/buddy/pkg/utils"
)

type ProjectConfig struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Scripts     map[string]string `json:"scripts"`
}

func NewProjectConfig(name string, version string, description string, author string, scripts map[string]string) *ProjectConfig {
	return &ProjectConfig{
		Name:        name,
		Version:     version,
		Description: description,
		Author:      author,
		Scripts:     scripts,
	}
}

func ParseProjectConfigFile(filePath string) (*ProjectConfig, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	content, err := io.ReadAll(readFile)
	if err != nil {
		return nil, err
	}

	var buddyConfig ProjectConfig
	if err := json.Unmarshal(content, &buddyConfig); err != nil {
		return nil, err
	}

	return &buddyConfig, nil
}

func (projectConfig *ProjectConfig) RunScript(scriptName string) error {
	return projectConfig.RunScriptArgs(scriptName, []string{})
}

func (projectConfig *ProjectConfig) RunScriptArgs(scriptName string, arguments []string) error {
	command, ok := projectConfig.Scripts[scriptName]
	if !ok {
		return fmt.Errorf("Script %s not found", scriptName)
	}

	_, err := utils.RunScript(command, arguments)

	return err
}

func (projectConfig *ProjectConfig) WriteToFile(filePath string) error {
	return utils.WriteJsonToFile(filePath, projectConfig)
}
