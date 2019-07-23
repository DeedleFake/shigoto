package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
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
	var dtype, title string
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

	buildPath, ok := t.get("buildPath").(string)
	if !ok {
		return errors.New("buildPath is not a string")
	}

	path, err := metaTmpl(buildPath, map[string]interface{}{
		"Type":  dtype,
		"Title": title,
		"Tmpl":  t.meta,
	})
	if err != nil {
		return fmt.Errorf("failed to construct buildPath: %v", err)
	}
	path = filepath.FromSlash(path)

	infile := filepath.Join(root, "draft", name)
	outfile := filepath.Join(root, "publish", filepath.Dir(path), name)

	in, err := os.Open(infile)
	if err != nil {
		return fmt.Errorf("failed to open %q: %v", name, err)
	}
	defer in.Close()

	meta := make(map[string]interface{}, 3)
	inr, err := readMeta(in, &meta)
	if err != nil {
		return fmt.Errorf("failed to read metadata from %q: %v", name, err)
	}
	if _, ok := meta["type"]; !ok {
		meta["type"] = dtype
	}
	if _, ok := meta["title"]; !ok {
		meta["title"] = title
	}
	if _, ok := meta["time"]; !ok {
		meta["time"] = time.Now().Format(time.RFC1123)
	}

	err = os.MkdirAll(filepath.Dir(outfile), 0755)
	if err != nil {
		return fmt.Errorf("failed to create %q: %v", path, err)
	}

	out, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create %q: %v", filepath.Join(filepath.Dir(path), name), err)
	}
	defer out.Close()

	e := yaml.NewEncoder(out)
	err = e.Encode(meta)
	if err != nil {
		return fmt.Errorf("failed to encode metadata: %v", err)
	}
	err = e.Close()
	if err != nil {
		return fmt.Errorf("failed to encode metadata: %v", err)
	}

	_, err = io.WriteString(out, "\n++++++++++\n")
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	_, err = io.Copy(out, inr)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	err = os.Remove(infile)
	if err != nil {
		return fmt.Errorf("failed to remove draft: %v", err)
	}

	return nil
}
