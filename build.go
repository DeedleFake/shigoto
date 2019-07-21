package main

import (
	"flag"
	"fmt"
)

type buildCmd struct {
	output string
}

func (cmd *buildCmd) Name() string {
	return "build"
}

func (cmd *buildCmd) Desc() string {
	return "generate static site"
}

func (cmd *buildCmd) Help() string {
	return `Usage: build [flags]:

	Builds the site, generating static pages to be served.`
}

func (cmd *buildCmd) Flags(fset *flag.FlagSet) {
	fset.StringVar(&cmd.output, "o", "build", "output directory")
}

func (cmd *buildCmd) Run(args []string) error {
	tmpl, err := loadTmpl()
	if err != nil {
		return fmt.Errorf("failed to load templates: %v", err)
	}

	fmt.Printf("%#v\n", tmpl)
	return nil

	//return build(tmpl)
}
