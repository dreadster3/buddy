package get

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthor(t *testing.T) {
	stdOutBuffer := &bytes.Buffer{}

	stdOutWriter := io.Writer(stdOutBuffer)

	expectedAuthor := "dreadster3"

	opts := &GetOptions{
		Settings: &settings.Settings{
			ProjectConfig: config.NewProjectConfig("buddy-tests", "0.0.1", "Tests for buddy", expectedAuthor, map[string]string{}),
			Logger:        slog.Default(),

			StdOut: stdOutWriter,
		},

		ParameterName: "author",
	}

	RunGet(opts)

	expected := fmt.Sprintf("%s\n", expectedAuthor)
	assert.Equal(t, expected, stdOutBuffer.String())
}

func TestGetScripts(t *testing.T) {
	stdOutBuffer := &bytes.Buffer{}

	stdOutWriter := io.Writer(stdOutBuffer)

	expectedAuthor := "dreadster3"

	opts := &GetOptions{
		Settings: &settings.Settings{
			ProjectConfig: config.NewProjectConfig("buddy-tests", "0.0.1", "Tests for buddy", expectedAuthor, map[string]string{}),
			Logger:        slog.Default(),

			StdOut: stdOutWriter,
		},

		ParameterName: "scripts",
	}

	err := RunGet(opts)

	if assert.Error(t, err) {
		assert.Equal(t, "Field is not printable", err.Error())
	}
}
