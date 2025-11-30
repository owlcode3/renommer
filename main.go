package main

import (
	"github.com/owlcode3/renommer/cmd"
	"github.com/owlcode3/renommer/internal/app"
	"github.com/owlcode3/renommer/internal/basepath"
	"github.com/spf13/afero"
)

func main() {
	app.Init(afero.NewOsFs(), basepath.GetWorkingDir())
	cmd.Execute()
}
