package core

import (
	"fmt"
	"os"
	"path/filepath"
)


const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite

	ColorBold     = 1
	ColorDarkGray = 90
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

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func Colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
