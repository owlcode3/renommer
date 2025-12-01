package pathx

import (
	"fmt"

	"github.com/owlcode3/renommer/internal/app"
	"github.com/owlcode3/renommer/internal/fs"
)

func RoutePath(args []string) error {
	_, err := app.A.FS.Stat(args[0])
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	err = fs.RenameEntry(args)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}
