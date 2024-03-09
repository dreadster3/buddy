package main

import (
	"runtime/debug"
)

// File to be used to automate versioning of the CLI tool
// This file will be updated by the build process to include the Version
var Version = "DEV"

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
