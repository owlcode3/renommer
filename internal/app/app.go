package app

import "github.com/spf13/afero"

type App struct {
	FS  afero.Fs
	CWD string
}

var A *App

func Init(fs afero.Fs, cwd string) {
	A = &App{
		FS:  fs,
		CWD: cwd,
	}
}
