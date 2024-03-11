package run

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestListAllCommands(t *testing.T) {
	stdOutBuffer := &bytes.Buffer{}

	stdOutWriter := io.Writer(stdOutBuffer)

	opts := &RunOptions{
		Settings: &settings.Settings{
			ProjectConfig: &config.ProjectConfig{
				Scripts: map[string]string{
					"script1": "echo 'script1'",
					"script2": "echo 'script2'",
				},
			},

			Logger: slog.Default(),
			StdOut: stdOutWriter,
		},

		ListCommands: true,
	}

	RunExecute(opts)

	expected := "script1  ->  echo 'script1'\nscript2  ->  echo 'script2'\n"

	assert.Equal(t, expected, stdOutBuffer.String())
}

func TestRunNoArgs(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := &RunOptions{
		Settings: &settings.Settings{
			ProjectConfig: &config.ProjectConfig{
				Scripts: map[string]string{
					"script1": "echo 'script1'",
				},
			},

			Logger: slog.Default(),
			StdOut: os.Stdout,
		},

		ScriptName: "script1",
		ScriptArgs: []string{},
	}

	RunExecute(opts)

	// Read the output
	w.Close()
	actual, _ := io.ReadAll(r)

	expected := "script1\n"

	assert.Equal(t, expected, string(actual))
}

func TestRunArgs(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := &RunOptions{
		Settings: &settings.Settings{
			ProjectConfig: &config.ProjectConfig{
				Scripts: map[string]string{
					"script1": "echo",
				},
			},

			Logger: slog.Default(),
			StdOut: os.Stdout,
		},

		ScriptName: "script1",
		ScriptArgs: []string{"hello", "world"},
	}

	RunExecute(opts)

	// Read the output
	w.Close()
	actual, _ := io.ReadAll(r)

	expected := "hello world\n"

	assert.Equal(t, expected, string(actual))
}
