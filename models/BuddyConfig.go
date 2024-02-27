package models

import (
	"encoding/json"
	"io"
	"os"
)

type BuddyConfig struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Scripts     map[string]string `json:"scripts"`
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
