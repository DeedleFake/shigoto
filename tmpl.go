package shigoto

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

	"github.com/DeedleFake/shigoto/internal/common"
	"github.com/gosimple/slug"
	"github.com/russross/blackfriday/v2"
)

func StandardFuncs(tmpls map[string]Tmpl) template.FuncMap {
	return template.FuncMap{
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

		"trimExt": func(file string) string {
			return strings.TrimSuffix(file, filepath.Ext(file))
		},

		"ext": func(file string) string {
			return strings.TrimPrefix(filepath.Ext(file), ".")
		},

		"tmpl": func(name string, data interface{}) (string, error) {
			if tmpls == nil {
				return "", errors.New("tmpl is unavailable in this context")
			}

			t, ok := tmpls[name]
			if !ok {
				return "", fmt.Errorf("unknown tmpl %q", name)
			}

			var out strings.Builder
			err := t.Tmpl.Execute(&out, data)
			return out.String(), err
		},

		// TODO: getByType(name string) ([]Content, error)
		// TODO: filter(k, check, val string, c []Content) []Content
		// TODO: slice(start, end int, c []Content) []Content
		// TODO: pageSlice(start, end int, c []Content) []Content
	}
}

type Tmpl struct {
	Meta map[string]interface{}
	Tmpl *template.Template
}

func LoadTmpl(root string) (map[string]Tmpl, error) {
	tmpls := make(map[string]Tmpl)
	err := common.Walk(root, func(path string, fi os.FileInfo) error {
		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(filepath.Join(root, path))
		if err != nil {
			return fmt.Errorf("failed to open %q: %v", path, err)
		}
		defer f.Close()

		var t Tmpl
		rem, err := ReadMeta(f, &t.Meta)
		if err != nil {
			return fmt.Errorf("failed to read meta from %q: %v", path, err)
		}

		var buf strings.Builder
		_, err = io.Copy(&buf, rem)
		if err != nil {
			return fmt.Errorf("failed to read %q: %v\n", path, err)
		}

		t.Tmpl = template.New(path)
		t.Tmpl.Funcs(StandardFuncs(tmpls))

		t.Tmpl, err = t.Tmpl.Parse(buf.String())
		if err != nil {
			return fmt.Errorf("failed to parse %q: %v", path, err)
		}

		tmpls[path] = t
		return nil
	})
	return tmpls, err
}
