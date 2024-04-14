package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"os/exec"

	"github.com/creack/pty"
	"github.com/dreadster3/buddy/pkg/log"
)

func RunScript(script string, arguments []string) (string, error) {
	log.Logger.Info("Running script", "script", script, "arguments", arguments)

	toRun := script
	if len(arguments) > 0 {
		toRun = fmt.Sprintf("%s %#v", script, strings.Join(arguments, " "))
	}

	command := exec.Command("bash", "-c", toRun)
	f, err := pty.Start(command)
	if err != nil {
		return "", err
	}

	io.Copy(os.Stdout, f)

	return string(""), nil
}
