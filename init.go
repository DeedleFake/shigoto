package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type initCmd struct {
	title string
}

func (cmd *initCmd) Name() string {
	return "init"
}

func (cmd *initCmd) Desc() string {
	return "creates a basic project structure"
}

func (cmd *initCmd) Help() string {
	return `Usage: init

The init command creates a basic project structure. It creates the
necessary directories, as well as a few basic example files.
`
}

func (cmd *initCmd) Flags(fset *flag.FlagSet) {
	fset.StringVar(&cmd.title, "title", "Example", "the default title for the site")
}

func (cmd *initCmd) Run(args []string) error {
	root := globalOptions.root

	files := []struct {
		path    string
		dir     bool
		content string
	}{
		{
			path: "tmpl",
			dir:  true,
		},
		{
			path: "draft",
			dir:  true,
		},
		{
			path: "publish",
			dir:  true,
		},

		{
			path: "tmpl/index.html",
			content: `<!DOCTYPE html>
<html>
	<head>
		<title>` + cmd.title + `</title>
	</head>
	<body>
		{{.Content | markdown}}
	</body>
</html>`,
		},
	}

	for _, file := range files {
		path := filepath.Join(root, filepath.FromSlash(file.path))

		if file.dir {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return fmt.Errorf("failed to create %q: %v", file.path, err)
			}
			continue
		}

		f, err := os.OpenFile(
			path,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0644,
		)
		if err != nil {
			return fmt.Errorf("failed to create %q: %v", file.path, err)
		}
		defer f.Close()

		_, err = io.Copy(f, strings.NewReader(file.content))
		if err != nil {
			return fmt.Errorf("failed to copy contents into %q: %v", file.path, err)
		}
	}

	return nil
}
