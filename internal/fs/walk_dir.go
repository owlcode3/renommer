package fs

import (
	"os"

	"github.com/owlcode3/renommer/internal/app"
	"github.com/spf13/afero"
)

type Entry struct {
	Name string
	Path string
	kind string
}

func Walk() ([]Entry, error) {
	var allEntry []Entry

	err := afero.Walk(
		app.A.FS,
		app.A.CWD,
		func(path string, d os.FileInfo, err error) error {
			if d.IsDir() {
				allEntry = append(allEntry, Entry{Name: d.Name(), Path: path, kind: "directory"})
			} else {
				allEntry = append(allEntry, Entry{Name: d.Name(), Path: path, kind: "file"})
			}

			if err != nil {
				return err
			}
			return nil
		})

	return allEntry, err
}
