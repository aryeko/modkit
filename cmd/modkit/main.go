// Package main is the entry point for the modkit CLI.
package main

import (
	"os"

	"github.com/go-modkit/modkit/internal/cli/cmd"
)

var osExit = os.Exit

func main() {
	osExit(run(os.Args))
}

func run(args []string) int {
	orig := os.Args
	os.Args = args
	defer func() {
		os.Args = orig
	}()

	if err := cmd.Execute(); err != nil {
		return 1
	}
	return 0
}
