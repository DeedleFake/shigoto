package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type buildCmd struct {
	output string
}

func (cmd *buildCmd) Name() string {
	return "build"
}

func (cmd *buildCmd) Desc() string {
	return "builds output"
}

func (cmd *buildCmd) Help() string {
	return `Usage: build [flags]:

The build command converts the content files in the publish directory
into static output files using the transformations specified in the
tmpl directory.`
}

func (cmd *buildCmd) Flags(fset *flag.FlagSet) {
	fset.StringVar(&cmd.output, "o", "build", "output directory name relative to project root")
}

func (cmd *buildCmd) Run(args []string) error {
	root, ok := getRoot()
	if !ok {
		return noRootErr
	}

	publish := filepath.Join(root, "publish")
	output := filepath.Join(root, cmd.output)

	tmpl, err := loadTmpl(root)
	if err != nil {
		return fmt.Errorf("failed to load templates: %v", err)
	}

	err = copyStatic(output, filepath.Join(root, "static"))
	if err != nil {
		return err
	}

	return walk(publish, func(p string, fi os.FileInfo) error {
		if fi.IsDir() {
			return nil
		}

		in, err := os.Open(filepath.Join(publish, p))
		if err != nil {
			return fmt.Errorf("failed to open %q: %v", p, err)
		}
		defer in.Close()

		meta := make(map[string]interface{})
		inr, err := readMeta(in, &meta)
		if err != nil {
			return fmt.Errorf("failed to load meta from %q: %v", p, err)
		}

		var instr strings.Builder
		_, err = io.Copy(&instr, inr)
		if err != nil {
			return fmt.Errorf("failed to read %q: %v", p, err)
		}

		intmpl, err := template.New(p).Funcs(standardFuncs(tmpl)).Parse(instr.String())
		if err != nil {
			return fmt.Errorf("failed to parse %q: %v", p, err)
		}

		dtype, ok := meta["type"].(string)
		if !ok {
			return fmt.Errorf("no type in %q", p)
		}

		title, _ := meta["title"].(string)

		t, ok := tmpl[dtype]
		if !ok {
			return fmt.Errorf("unknown type %q in %q", dtype, p)
		}

		pages, ok := tmplGet("pages", meta, t.meta).(pagesInfo)
		if !ok {
			return fmt.Errorf("pages is not an object in %q", p)
		}

		numPages := 1

		var numType int
		if pages.Tmpl != "" {
			num, err := getNumType(pages.Tmpl)
			if err != nil {
				return fmt.Errorf("failed to get number of pages for %q", pages.Tmpl)
			}
			numType = num

			var extra int
			if numType%pages.Per != 0 {
				extra = 1
			}

			numPages = (numType / pages.Per) + extra
		}

		for currentPage := 1; currentPage <= numPages; currentPage++ {
			pageEnd := currentPage * pages.Per
			if pageEnd > numType {
				pageEnd = numType
			}

			pageMap := map[string]interface{}{
				"Last":      numPages,
				"Current":   currentPage,
				"PageStart": (currentPage - 1) * pages.Per,
				"PageEnd":   pageEnd,
			}

			var content strings.Builder
			err = intmpl.Execute(&content, map[string]interface{}{
				"Type":  dtype,
				"Title": title,
				"Tmpl":  t.meta,
				"Meta":  meta,
				"Pages": pageMap,
			})
			if err != nil {
				return fmt.Errorf("failed to execute %q: %v", dtype, err)
			}

			buildPath, ok := tmplGet("buildPath", meta, t.meta).(string)
			if !ok {
				return fmt.Errorf("buildPath is not a string in %q", p)
			}

			path, err := metaTmpl(buildPath, map[string]interface{}{
				"Type":  dtype,
				"Title": title,
				"Tmpl":  t.meta,
				"Meta":  meta,
				"Pages": pageMap,
			})
			if err != nil {
				return fmt.Errorf("failed to construct buildPath for %q: %v", p, err)
			}
			path = filepath.FromSlash(path)

			err = os.MkdirAll(filepath.Join(output, filepath.Dir(path)), 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory for %q: %v", p, err)
			}

			out, err := os.OpenFile(
				filepath.Join(output, path),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				0644,
			)
			if err != nil {
				return fmt.Errorf("failed to create %q: %v", p, err)
			}
			defer out.Close()

			err = executeInherit(tmpl, t, out, map[string]interface{}{
				"Type":    dtype,
				"Title":   title,
				"Tmpl":    t.meta,
				"Meta":    meta,
				"Content": content.String(),
				"Pages":   pageMap,
			})
			if err != nil {
				return fmt.Errorf("failed to execute %q: %v", p, err)
			}
		}

		return nil
	})
}

func copyStatic(out, in string) error {
	_, err := os.Stat(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "No static directory found.")
		return nil
	}

	return walk(in, func(p string, fi os.FileInfo) error {
		err := os.MkdirAll(filepath.Join(out, filepath.Dir(p)), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory for %q: %v", p, err)
		}

		err = os.RemoveAll(filepath.Join(out, p))
		if err != nil {
			return fmt.Errorf("failed to remove %q: %v", p, err)
		}

		err = os.Link(filepath.Join(in, p), filepath.Join(out, p))
		if err != nil {
			return fmt.Errorf("failed to copy %q: %v", p, err)
		}

		return nil
	})
}

func executeInherit(tmpl map[string]tmpl, t tmpl, out io.Writer, data map[string]interface{}) error {
	inherit, ok := tmplGet("inherit", t.meta).(string)
	if !ok {
		return t.tmpl.Execute(out, data)
	}

	next, ok := tmpl[inherit]
	if !ok {
		return fmt.Errorf("unknown type %q", inherit)
	}

	var content strings.Builder
	err := t.tmpl.Execute(&content, data)
	if err != nil {
		return err
	}

	nextData := make(map[string]interface{}, len(data))
	for k, v := range data {
		nextData[k] = v
	}
	for k, v := range next.meta {
		nextData["Tmpl"].(map[string]interface{})[k] = v
	}
	nextData["Content"] = content.String()

	return executeInherit(tmpl, next, out, nextData)
}
