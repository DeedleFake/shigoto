package main

import (
	"flag"
	"fmt"
	"os"
)

type publishCmd struct{}

func (cmd *publishCmd) Name() string {
	return "publish"
}

func (cmd *publishCmd) Desc() string {
	return "publishes a draft"
}

func (cmd *publishCmd) Help() string {
	return `Usage: publish <type> <title>

The publish command publishes an existing draft by moving it from the
draft directory to the publish directory. It puts it into a directory
that matches where its output will be placed in the build directory
when the project is built. It also inserts a timestamp into the
draft's metadata with the name "time", unless an entry in the metadata
with that name already exists.`
}

func (cmd *publishCmd) Flags(fset *flag.FlagSet) {
	// TODO: -flat for not putting published drafts into subdirectories.
	// TODO: -time=false for not inserting the publish timestamp.
}

func (cmd *publishCmd) Run(args []string) error {
	switch len(args) {
	case 0, 1:
		fmt.Fprintf(os.Stderr, "Error: not enough arguments\n\n")
		return flag.ErrHelp

	case 2:
		dtype = args[0]
		title = args[1]

	default:
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n\n")
		return flag.ErrHelp
	}

	return nil
}
