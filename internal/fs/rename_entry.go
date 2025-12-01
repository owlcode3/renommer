package fs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/owlcode3/renommer/internal/app"
	"github.com/owlcode3/renommer/internal/common"
)

func RenameEntry(args []string) error {
	var oldPath, newname = args[0], args[1]

	if !filepath.IsAbs(oldPath) {
		oldPath = filepath.Join(app.A.CWD, oldPath)
	}

	oldPath = filepath.Clean(oldPath)

	if filepath.Clean(oldPath) == app.A.CWD {
		return fmt.Errorf("cannot rename project root directory")
	}

	var oldPathExist bool
	var entryKind string

	allEntry, err := Walk()
	if err != nil {
		return err
	}

	for _, entry := range allEntry {
		if oldPath == entry.Path {
			oldPathExist = true
			entryKind = entry.kind
			break
		}
	}

	if !oldPathExist {
		return fmt.Errorf("unable to find '%s' in current working directory", oldPath)
	}

	var newPath string

	if strings.Contains(newname, string(filepath.Separator)) &&
		filepath.Dir(oldPath) == filepath.Dir(filepath.Clean(newname)) {
		newPath = newname
	} else if !strings.Contains(newname, string(filepath.Separator)) && strings.TrimSpace(newname) != "" {
		newPath = filepath.Join(filepath.Dir(oldPath), newname)
	} else {
		return fmt.Errorf("unable to rename directory/file to '%s' as name pattern is unsupported", newname)
	}

	for _, entry := range allEntry {
		if entry.Path == filepath.Join(filepath.Dir(oldPath), filepath.Base(newPath)) {
			return fmt.Errorf("a %s '%s' already exist at this location. please choose a different name", entry.kind, newPath)
		}
	}

	relOldPath, _ := filepath.Rel(app.A.CWD, oldPath)
	relNewPath, _ := filepath.Rel(app.A.CWD, newPath)

	if common.Verbose {
		color.Cyan("Renaming %s '%v' to '%v' .......\n", entryKind, relOldPath, relNewPath)
	}

	if err := app.A.FS.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("error renaming '%s' '%s' to '%s': %v", entryKind, relOldPath, relNewPath, err)
	}

	entryKind = strings.ToUpper(string(entryKind[0])) + entryKind[1:]

	if common.Verbose {
		color.Green("%s '%s' successfully renamed to '%s'", entryKind, relOldPath, relNewPath)
	}

	return nil
}
