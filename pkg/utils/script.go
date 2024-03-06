package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func RunScript(script string, arguments []string) (string, error) {
	var stdOutBuffer, stdErrBuffer bytes.Buffer

	stdOutWriter := io.MultiWriter(os.Stdout, &stdOutBuffer)
	stdErrWriter := io.MultiWriter(os.Stderr, &stdErrBuffer)

	script = fmt.Sprintf("%s %s", script, strings.Join(arguments, " "))

	execCommand := exec.Command("sh", "-c", script)
	execCommand.Stdout = stdOutWriter
	execCommand.Stderr = stdErrWriter

	err := execCommand.Run()
	if err != nil {
		return "", errors.New(stdErrBuffer.String())
	}

	return stdOutBuffer.String(), nil
}
