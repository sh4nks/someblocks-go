package core

import (
	"os"
	"path/filepath"
)

func GetAppDir() string {
    var dirAbsPath string
    ex, err := os.Executable()
    if err == nil {
        dirAbsPath = filepath.Dir(ex)
        return dirAbsPath
    }

    exReal, err := filepath.EvalSymlinks(ex)
    if err != nil {
        panic(err)
    }
    dirAbsPath = filepath.Dir(exReal)
    return dirAbsPath
}
