package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DeedleFake/sub"
)

var (
	noRootErr = errors.New("couldn't find root of project")
)

var globalOptions struct {
	root string
}

func main() {
	commander := &sub.Commander{
		Help: `
shigoto is a simple static site generator designed primarily for when
Hugo is just complete overkill. It has no config files, instead
relying entirely on configs embedded into the source files themselves
or other, similar sources.
	`,

		Flags: func(fset *flag.FlagSet) {
			fset.StringVar(&globalOptions.root, "root", "", "the root of the project")
		},
	}

	commander.Register(commander.HelpCmd())
	commander.Register(&initCmd{})
	commander.Register(&draftCmd{})
	commander.Register(&publishCmd{})
	commander.Register(&buildCmd{})
	commander.Register(&cleanCmd{})

	err := commander.Run(append([]string{filepath.Base(os.Args[0])}, os.Args[1:]...))
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(2)
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
