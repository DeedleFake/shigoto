package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type draftCmd struct{}

func (cmd *draftCmd) Name() string {
	return "draft"
}

func (cmd *draftCmd) Desc() string {
	return "creates a new draft"
}

func (cmd *draftCmd) Help() string {
	return `Usage: draft <type> [title]

The draft command creates a new draft of the given type and prints the
path to it to stdout. If a file already exists at the path that the
new draft would be created at, the path to that file is printed but no
other action is taken.

For example, if you want to edit, using vim, a draft of type page.html
with the title "This is an Example", regardless of whether it exists
or not, simply run

    $ vim $(shigoto draft "This is an Example")`
}

func (cmd *draftCmd) Flags(fset *flag.FlagSet) {
}

func (cmd *draftCmd) Run(args []string) error {
	var dtype, title string
	switch len(args) {
	case 0:
		fmt.Fprintf(os.Stderr, "Error: must specify draft type\n\n")
		return flag.ErrHelp

	case 1:
		dtype = args[0]
		title = time.Now().Format("2006-01-02-15-04")

	case 2:
		dtype = args[0]
		title = args[1]

	default:
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n\n")
		return flag.ErrHelp
	}

	root, ok := getRoot()
	if !ok {
		return noRootErr
	}

	tmpl, err := loadTmpl(root)
	if err != nil {
		return err
	}

	t, ok := tmpl[dtype]
	if !ok {
		return fmt.Errorf("unknown type %q", dtype)
	}

	sourceName, ok := t.get("sourceName").(string)
	if !ok {
		return errors.New("sourceName is not a string")
	}

	name, err := metaTmpl(sourceName, map[string]interface{}{
		"Type":  dtype,
		"Title": title,
		"Tmpl":  t.meta,
	})
	if err != nil {
		return err
	}

	path := filepath.Join(root, "draft", name)

	_, err = os.Stat(path)
	if err == nil {
		fmt.Println(path)
		return nil
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create draft: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(
		file,
		`type: %q
title: %q

++++++++++`,
		dtype,
		title,
	)
	if err != nil {
		return fmt.Errorf("failed to write draft skeleton: %v", err)
	}

	fmt.Println(path)

	return nil
}
