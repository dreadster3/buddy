package main

import (
	"os"

	"github.com/dreadster3/buddy/pkg/cmd/get"
	"github.com/dreadster3/buddy/pkg/cmd/initialize"
	"github.com/dreadster3/buddy/pkg/cmd/root"
	"github.com/dreadster3/buddy/pkg/cmd/run"
)

func main() {
	cmd := root.NewRootCmd()

	cmd.AddCommand(initialize.NewCmdInit())
	cmd.AddCommand(get.NewCmdGet())
	cmd.AddCommand(run.NewCmdRun())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
