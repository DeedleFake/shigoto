package shigoto

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/DeedleFake/shigoto/internal/common"
	"gopkg.in/yaml.v2"
)

var metaSplit = regexp.MustCompile(`^\+{5,}\n?$`)

func ReadMeta(r io.Reader, v interface{}) (rem io.Reader, err error) {
	br := bufio.NewReader(r)

	var buf []byte
	for {
		line, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if metaSplit.Match(line) {
					break
				}

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

func FindRoot(path string) (string, bool) {
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

var numTypeCache struct {
	common.Once
	c map[string]int
}

func GetNumType(root string, name string) (int, error) {
	err := numTypeCache.Do(func() error {
		counts := make(map[string]int)
		err := common.Walk(filepath.Join(root, "publish"), func(p string, fi os.FileInfo) error {
			if fi.IsDir() {
				return nil
			}

			f, err := os.Open(filepath.Join(root, "publish", p))
			if err != nil {
				return fmt.Errorf("failed to open %q: %v", p, err)
			}
			defer f.Close()

			var meta map[string]interface{}
			_, err = ReadMeta(f, &meta)
			if err != nil {
				return fmt.Errorf("failed to read metadata from %q: %v", p, err)
			}

			dtype, _ := meta["type"].(string)
			counts[dtype]++

			return nil
		})
		if err != nil {
			return err
		}

		numTypeCache.c = counts
		return nil
	})
	if err != nil {
		return 0, err
	}

	return numTypeCache.c[name], nil
}
