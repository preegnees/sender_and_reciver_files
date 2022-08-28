package common

import (
	"strings"
)

func GetRelativePath(parentDir, pathToFile string) (relativePath string) {
	return strings.Replace(pathToFile, parentDir, "", 1)
}