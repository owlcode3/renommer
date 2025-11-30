package fs_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/owlcode3/renommer/internal/app"
	"github.com/owlcode3/renommer/internal/fs"
	"github.com/spf13/afero"
)

func TestRenameEntry(t *testing.T) {
	testCases := []struct {
		name    string
		old     string
		new     string
		wantErr bool
		setup   func(FS afero.Fs)
		verify  func(FS afero.Fs) error
	}{
		{
			name:    "rename directory with nested items",
			old:     "/project/test",
			new:     "testing",
			wantErr: false,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project/test/this/that", 0755)
			},
			verify: func(FS afero.Fs) error {
				if _, err := FS.Stat("/project/testing"); err != nil {
					return fmt.Errorf("expected to exist after rename")
				}

				if _, err := FS.Stat("/project/test"); !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("expected to not exist after rename")
				}

				if _, err := FS.Stat("/project/testing/this/that"); err != nil {
					return fmt.Errorf("expected nested structure preserved: %v", err)
				}

				return nil
			},
		},
		{
			name:    "rename file",
			old:     "/project/test/file.txt",
			new:     "/project/test/renamed.txt",
			wantErr: false,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project/test", 0755)
				file, _ := FS.Create("/project/test/file.txt")
				defer file.Close()
			},
			verify: func(FS afero.Fs) error {
				if _, err := FS.Stat("/project/test/renamed.txt"); err != nil {
					return fmt.Errorf("expected file renamed: %v", err)
				}

				if _, err := FS.Stat("/project/test/file.txt"); !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("old file still exists")
				}

				return nil
			},
		},
		{
			name:    "rename project root directory",
			old:     "/project",
			new:     "something",
			wantErr: true,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project", 0755)
			},
			verify: nil,
		},
		{
			name:    "rename directory using absolute dest path",
			old:     "/project/test",
			new:     "/project/testing",
			wantErr: false,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project/test/n1/n2", 0755)
			},
			verify: func(FS afero.Fs) error {
				if _, err := FS.Stat("/project/testing/n1/n2"); err != nil {
					return fmt.Errorf("absolute rename failed: %v", err)
				}
				return nil
			},
		},
		{
			name:    "non-existent source",
			old:     "/project/nope",
			new:     "testing",
			wantErr: true,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project", 0755)
			},
			verify: nil,
		},
		{
			name:    "dest already exists",
			old:     "/project/yikes",
			new:     "jeez",
			wantErr: true,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project/yikes", 0755)
				FS.MkdirAll("/project/jeez", 0755)
			},
			verify: nil,
		},
		{
			name:    "empty args",
			old:     "",
			new:     "",
			wantErr: true,
			setup: func(FS afero.Fs) {
				FS.MkdirAll("/project", 0755)
			},
			verify: nil,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Init(afero.NewMemMapFs(), "/project")

			FS := app.A.FS
			tc.setup(FS)

			err := fs.RenameEntry([]string{tc.old, tc.new})
			if tc.wantErr && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("did not expect error but got: %v", err)
			}

			if tc.verify != nil {
				if err := tc.verify(FS); err != nil {
					t.Fatalf("verification failed: %v", err)
				}
			}
		})
	}
}
