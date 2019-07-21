package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unsafe"

	"github.com/russross/blackfriday"
)

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

		var buf strings.Builder
		_, err = io.Copy(&buf, rem)
		if err != nil {
			return fmt.Errorf("failed to read %q: %v\n", path, err)
		}

		t.tmpl = template.New(path)
		t.tmpl.Funcs(standardFuncs)

		t.tmpl, err = t.tmpl.Parse(buf.String())
		if err != nil {
			return fmt.Errorf("failed to parse %q: %v", path, err)
		}

		tmpls[path] = t
		return nil
	})
	return tmpls, err
}

var standardFuncs = map[string]interface{}{
	"markdown": func(str string) string {
		out := blackfriday.Run([]byte(str))
		return *(*string)(unsafe.Pointer(&out))
	},
}
