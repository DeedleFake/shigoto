package common

import (
	"os"
	"path/filepath"
)

func Walk(root string, f func(path string, fi os.FileInfo) error) error {
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
