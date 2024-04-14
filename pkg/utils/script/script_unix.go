//go:build !windows

package script

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"os/exec"

	"github.com/creack/pty"
	"github.com/dreadster3/buddy/pkg/log"
)

func RunScript(script string, arguments []string) (string, string, error) {
	var stdOutBuffer, stdErrBuffer bytes.Buffer

	ptyStdOut, ttyStdOut, err := pty.Open()
	if err != nil {
		return "", "", err
	}

	ptyStdErr, ttyStdErr, err := pty.Open()
	if err != nil {
		return "", "", err
	}

	stdOutMultiWriter := io.MultiWriter(os.Stdout, &stdOutBuffer)
	stdErrMultiWriter := io.MultiWriter(os.Stderr, &stdErrBuffer)

	log.Logger.Info("Running script", "script", script, "arguments", arguments)

	toRun := script
	if len(arguments) > 0 {
		toRun = fmt.Sprintf("%s %#v", script, strings.Join(arguments, " "))
	}

	command := exec.Command("bash", "-c", toRun)
	command.Stdout = ttyStdOut
	command.Stderr = ttyStdErr

	go func() { io.Copy(stdOutMultiWriter, ptyStdOut) }()
	go func() { io.Copy(stdErrMultiWriter, ptyStdErr) }()
	err = command.Run()
	if err != nil {
		return stdOutBuffer.String(), stdErrBuffer.String(), nil
	}

	return stdOutBuffer.String(), stdErrBuffer.String(), nil
}
