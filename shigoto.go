package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DeedleFake/sub"
)

func main() {
	var commander sub.Commander
	commander.Help(`
shigoto is a simple static site generator designed primarily for when
Hugo is just complete overkill. It has no config files, instead
relying entirely on configs embedded into the source files themselves
or other, similar sources.
`)

	commander.Register(commander.HelpCmd())
	commander.Register(&buildCmd{})

	err := commander.Run(append([]string{filepath.Base(os.Args[0])}, os.Args[1:]...))
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(2)
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
