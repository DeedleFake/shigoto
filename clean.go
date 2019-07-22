package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type cleanCmd struct {
	output string
}

func (cmd *cleanCmd) Name() string {
	return "clean"
}

func (cmd *cleanCmd) Desc() string {
	return "removes the build directory"
}

func (cmd *cleanCmd) Help() string {
	return `Usage: clean [options]

The clean command is purely for convience. It locates and removes the
project's output directory. In many situations, this may not be
particularly useful, but if you're several levels down into the
directory hierarchy of your project and you want to delete the output
directory, this will do so.

Warning: This command simply locates the project root and deletes the
specified directory relative to that. Make sure you don't tell it to
delete the wrong one by accident.`
}

func (cmd *cleanCmd) Flags(fset *flag.FlagSet) {
	fset.StringVar(&cmd.output, "o", "build", "directory to remove relative to root")
}

func (cmd *cleanCmd) Run(args []string) error {
	root, ok := getRoot()
	if !ok {
		return noRootErr
	}

	err := os.RemoveAll(filepath.Join(root, cmd.output))
	if err != nil {
		return fmt.Errorf("failed to remove %q: %v", cmd.output, err)
	}

	return nil
}
