package initialize

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	path, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	err = os.Chdir(path)
	if err != nil {
		t.Fatal(err)
	}

	stdOutBuffer := &bytes.Buffer{}
	stdOutWriter := io.Writer(stdOutBuffer)

	opts := &InitOptions{
		Settings: &settings.Settings{
			GlobalConfig: &config.GlobalConfig{
				Author:   "dreadster3",
				FileName: "buddy.json",
			},
			Logger:           slog.Default(),
			WorkingDirectory: ".",

			StdOut: stdOutWriter,
		},

		ProjectName: "buddy-tests",
		Description: "Description",
	}

	err = RunInit(opts)

	assert.Nil(t, err)
	assert.Contains(t, "buddy.json created\n", stdOutBuffer.String())

	file, err := os.Open(opts.Settings.GlobalConfig.FileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)

	var fileInterface map[string]interface{}

	err = json.Unmarshal(fileContent, &fileInterface)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "buddy-tests", fileInterface["name"])
	assert.Equal(t, "Description", fileInterface["description"])
	assert.Equal(t, "dreadster3", fileInterface["author"])
	assert.Equal(t, "0.0.1", fileInterface["version"])
}
