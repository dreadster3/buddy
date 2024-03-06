package utils

import (
	"encoding/json"
	"os"
)

func WriteJsonToFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	json, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(json))
	if err != nil {
		return err
	}

	return nil
}
