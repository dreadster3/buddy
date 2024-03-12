package root

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestRunRootNoArgs(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := &RootOptions{
		Settings: &settings.Settings{
			ProjectConfig: &config.ProjectConfig{
				Scripts: map[string]string{
					"script1": "echo 'hello world'",
				},
			},

			Logger: slog.Default(),
			StdOut: os.Stdout,
		},

		ScriptName: "script1",
		ScriptArgs: []string{},
	}

	err := RunRoot(opts)

	// Read the output
	w.Close()

	assert.Nil(t, err)

	actual, _ := io.ReadAll(r)
	expected := "hello world\n"

	assert.Equal(t, expected, string(actual))
}

func TestRunRootArgs(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := &RootOptions{
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
		ScriptArgs: []string{"Hello", "World"},
	}

	err := RunRoot(opts)

	// Read the output
	w.Close()

	assert.Nil(t, err)

	actual, _ := io.ReadAll(r)
	expected := "Hello World\n"

	assert.Equal(t, expected, string(actual))
}
