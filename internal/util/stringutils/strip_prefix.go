package stringutils

import (
	"os"
	"strings"
)

func StripPrefixDirPath(path string, dir string) string {
	prefix := dir + string(os.PathSeparator)

	if strings.HasPrefix(path, prefix) {
		return strings.TrimPrefix(path, prefix)
	}

	return path
}
