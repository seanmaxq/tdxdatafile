package tdxdatafile

import (
	"runtime"
)

var (
	PATH_SEPARATOR = GetPathSeparator()
)

func GetPathSeparator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}
