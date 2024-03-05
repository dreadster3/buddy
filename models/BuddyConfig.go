package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type BuddyConfig struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Scripts     map[string]string `json:"scripts"`
}

func NewBuddyConfig(name string, version string, description string, author string, scripts map[string]string) *BuddyConfig {
	return &BuddyConfig{
		Name:        name,
		Version:     version,
		Description: description,
		Author:      author,
		Scripts:     scripts,
	}
}

func ParseBuddyConfigFile(filePath string) (*BuddyConfig, error) {
	readFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	content, err := io.ReadAll(readFile)
	if err != nil {
		return nil, err
	}

	var buddyConfig BuddyConfig
	if err := json.Unmarshal(content, &buddyConfig); err != nil {
		return nil, err
	}

	return &buddyConfig, nil
}

func (buddyConfig *BuddyConfig) RunScript(scriptName string) error {
	return buddyConfig.RunScriptArgs(scriptName, []string{})
}

func (buddyConfig *BuddyConfig) RunScriptArgs(scriptName string, arguments []string) error {
	command, ok := buddyConfig.Scripts[scriptName]
	if !ok {
		return fmt.Errorf("Script %s not found", scriptName)
	}

	command = fmt.Sprintf("%s %s", command, strings.Join(arguments, " "))

	execCommand := exec.Command("sh", "-c", command)
	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr

	err := execCommand.Run()
	if err != nil {
		return err
	}

	return nil
}

func (buddyConfig *BuddyConfig) ToJson() ([]byte, error) {
	json, err := json.MarshalIndent(buddyConfig, "", "    ")
	if err != nil {
		return nil, err
	}

	return json, nil
}
