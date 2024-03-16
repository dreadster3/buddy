package initialize

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/dreadster3/buddy/pkg/cmd/settings"
	"github.com/dreadster3/buddy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func createTempFolder() (string, error) {
	path, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	err = os.Chdir(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func createTemplateFolder(templatePath string, templateName string) error {
	_ = os.MkdirAll(filepath.Join(templatePath, templateName), 0755)

	file, err := os.Create(filepath.Join(templatePath, templateName, "test"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("{{.ProjectConfig.Name}}")
	if err != nil {
		return err
	}

	return nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func TestInit(t *testing.T) {
	path, err := createTempFolder()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	stdOutBuffer := &bytes.Buffer{}
	stdOutWriter := io.Writer(stdOutBuffer)

	opts := &InitOptions{
		Settings: &settings.Settings{
			GlobalConfig: &config.GlobalConfig{
				Author:   "dreadster3",
				FileName: "buddy.json",
			},
			Logger:           slog.Default(),
			WorkingDirectory: "buddy-tests",

			StdOut: stdOutWriter,
		},

		ProjectName: "buddy-tests",
		Description: "Description",
	}

	os.MkdirAll(opts.Settings.WorkingDirectory, 0755)

	err = RunInit(opts)

	assert.Nil(t, err)
	assert.Contains(t, "Project initialized successfully!\n", stdOutBuffer.String())

	fileContent, err := readFile(filepath.Join(opts.Settings.WorkingDirectory, opts.Settings.GlobalConfig.FileName))

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

func TestTemplateInit(t *testing.T) {
	tempFolderPath, err := createTempFolder()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempFolderPath)

	t.Logf("tempFolderPath: %s", tempFolderPath)

	templateName := "test-template"

	stdOutBuffer := &bytes.Buffer{}
	stdOutWriter := io.Writer(stdOutBuffer)

	templatesPath, err := filepath.Abs(filepath.Join(tempFolderPath, "templates"))
	t.Logf("templatesPath: %s", templatesPath)

	opts := &InitOptions{
		Settings: &settings.Settings{
			GlobalConfig: &config.GlobalConfig{
				Author:        "dreadster3",
				TemplatesPath: templatesPath,
				FileName:      "buddy.json",
			},
			Logger:           slog.Default(),
			WorkingDirectory: filepath.Join(tempFolderPath, "buddy-tests"),
			StdOut:           stdOutWriter,
		},

		ProjectName:  "buddy-tests",
		Description:  "Description",
		TemplateName: templateName,
	}

	t.Logf("Opts: %+v", opts)
	os.MkdirAll(opts.Settings.WorkingDirectory, 0755)

	err = createTemplateFolder(templatesPath, templateName)
	if err != nil {
		t.Fatal(err)
	}

	err = RunInit(opts)

	file, err := os.Open(filepath.Join(opts.Settings.WorkingDirectory, opts.Settings.GlobalConfig.FileName))
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

	assert.Nil(t, err)
	assert.Equal(t, "buddy-tests", fileInterface["name"])
	assert.Equal(t, "Description", fileInterface["description"])
	assert.Equal(t, "dreadster3", fileInterface["author"])
	assert.Equal(t, "0.0.1", fileInterface["version"])

	renderedFile, err := readFile(filepath.Join(opts.Settings.WorkingDirectory, "test"))

	assert.Nil(t, err)
	assert.Equal(t, "buddy-tests", string(renderedFile))
}
