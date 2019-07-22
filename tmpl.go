package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/gosimple/slug"
	"github.com/russross/blackfriday"
)

var standardFuncs = map[string]interface{}{
	"markdown": func(str string) string {
		out := blackfriday.Run([]byte(str))
		return *(*string)(unsafe.Pointer(&out))
	},

	"slug": slug.Make,

	"time": func(t interface{}) (time.Time, error) {
		switch t := t.(type) {
		case int:
			return time.Unix(int64(t), 0), nil

		case string:
			for _, f := range []string{time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano} {
				t, err := time.Parse(f, t)
				if err != nil {
					continue
				}

				return t, nil
			}

			return time.Time{}, errors.New("failed to parse time")

		default:
			return time.Time{}, fmt.Errorf("unexpected time type: %T", t)
		}
	},
}

var defaults = map[string]interface{}{
	"sourceName": `{{.Title | slug}}.md`,
}

type tmpl struct {
	tmpl *template.Template
	meta map[string]interface{}
}

func loadTmpl(root string) (map[string]tmpl, error) {
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
		rem, err := readMeta(f, &t.meta)
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

func (t tmpl) get(name string) interface{} {
	sourceName, ok := t.meta[name]
	if !ok {
		return defaults[name]
	}
	return sourceName
}

func metaTmpl(src string, data interface{}) (string, error) {
	snt, err := template.New(src).Funcs(standardFuncs).Parse(src)
	if err != nil {
		return "", err
	}

	var r strings.Builder
	err = snt.Execute(&r, data)
	return r.String(), err
}
