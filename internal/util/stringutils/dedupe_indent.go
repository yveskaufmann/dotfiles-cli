package stringutils

import "strings"

// CommonIndent returns the number of leading spaces/tabs that are common across all non-empty lines in the input string.
func CommonIndent(str string) int {
	lines := strings.Split(str, "\n")
	minIndent := -1

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue // skip empty lines
		}
		indent := len(line) - len(trimmed)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}
	return minIndent
}

// DepdupeIndention removes the common leading indentation from all non-empty lines in the input string
// Helps to build multi-line strings without unwanted indention to the left.
func DepdupeIndention(str string) string {
	lines := strings.Split(str, "\n")
	var dedupedLines []string

	minIndent := CommonIndent(str)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if len(line) >= minIndent {
			dedupedLines = append(dedupedLines, line[minIndent:])
		}
	}

	return strings.Join(dedupedLines, "\n")
}
