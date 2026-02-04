package theme

import (
	"os"

	"golang.org/x/term"
)

// ANSI color codes
const (
	ColorReset  = "\x1b[0m"
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorCyan   = "\x1b[36m"
)

var UseColors = term.IsTerminal(int(os.Stdout.Fd()))

// Colorize returns the color code if colors are enabled, otherwise empty string
func Colorize(color string) string {
	if UseColors {
		return color
	}
	return ""
}
