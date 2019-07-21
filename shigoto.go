package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/DeedleFake/sub"
	"github.com/go-yaml/yaml"
)

var metaSplit = regexp.MustCompile(`^\+{5,}\n$`)

func readMeta(r io.ReadCloser, v interface{}) (rem io.Reader, err error) {
	br := bufio.NewReader(r)

	var buf []byte
	for {
		line, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				buf = append(buf, line...)
				return bytes.NewBuffer(buf), nil
			}

			return rem, err
		}

		if metaSplit.Match(line) {
			break
		}

		buf = append(buf, line...)
	}

	err = yaml.Unmarshal(buf, v)
	return br, err
}

type tmpl struct {
	tmpl *template.Template

	Path string `yaml:"path"`
}

func loadTmpl() (map[string]tmpl, error) {
	root, ok := getRoot()
	if !ok {
		return nil, errors.New("couldn't find root of project")
	}
	root = filepath.Join(root, "tmpl")

	tmpls := make(map[string]tmpl)
	err := walk(root, func(path string, fi os.FileInfo) error {
		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(filepath.Join(root, path))
		if err != nil {
			return fmt.Errorf("failed to open %q: %v", path, err)
		}
		defer f.Close()

		var t tmpl
		rem, err := readMeta(f, &t)
		if err != nil {
			return fmt.Errorf("failed to read meta from %q: %v", path, err)
		}
		tmpls[path] = t

		io.Copy(os.Stdout, rem)

		return nil
	})
	return tmpls, err
}

func getRoot() (string, bool) {
	path := globalOptions.root
	if path == "" {
		path, _ = os.Getwd()
	}
	path = filepath.Clean(path)

	for {
		d, err := os.Open(path)
		if err != nil {
			continue
		}
		defer d.Close()

		fi, err := d.Readdir(-1)
		if err != nil {
			continue
		}

		var found int
		for _, info := range fi {
			switch info.Name() {
			case "tmpl", "draft", "publish":
				if info.IsDir() {
					found++
					if found == 3 {
						return path, true
					}
				}
			}
		}

		next := filepath.Dir(path)
		if next == path {
			return "", false
		}
		path = next
	}
}

func walk(root string, f func(path string, fi os.FileInfo) error) error {
	var inner func(cur string) error
	inner = func(cur string) error {
		d, err := os.Open(cur)
		if err != nil {
			return err
		}
		defer d.Close()

		entries, err := d.Readdir(-1)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			p, _ := filepath.Rel(root, cur)
			err := f(filepath.Join(p, entry.Name()), entry)
			if err != nil {
				return err
			}

			if entry.IsDir() {
				err := inner(filepath.Join(cur, entry.Name()))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
	return inner(root)
}

var globalOptions struct {
	root string
}

func main() {
	var commander sub.Commander
	commander.Help(`
shigoto is a simple static site generator designed primarily for when
Hugo is just complete overkill. It has no config files, instead
relying entirely on configs embedded into the source files themselves
or other, similar sources.
	`)

	commander.Flags(func(fset *flag.FlagSet) {
		fset.StringVar(&globalOptions.root, "root", "", "the root of the project")
	})

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
