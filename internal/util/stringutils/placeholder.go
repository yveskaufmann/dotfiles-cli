package stringutils

import (
	"runtime"
	"strings"
)

// ResolvePlaceholders replaces placeholders like :arch, :os from the current system.
func ResolvePlaceholders(str string) string {
	str = strings.ReplaceAll(str, "{os}", runtime.GOOS)
	str = strings.ReplaceAll(str, "{arch}", runtime.GOARCH)
	return str
}

// ResolvePlaceholdersWithVars replaces placeholders like {arch}, {os} and any in vars.
func ResolvePlaceholdersWithVars(str string, vars map[string]string) string {
	str = ResolvePlaceholders(str)

	for k, v := range vars {
		k = strings.TrimPrefix(k, "{")
		k = strings.TrimSuffix(k, "}")
		str = strings.ReplaceAll(str, "{"+k+"}", v)
	}

	return str
}
