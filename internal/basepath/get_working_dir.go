package basepath

import (
	"os"
)

func GetWorkingDir() string {
	var cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	return cwd
}
