package main

import (
	"yv35.com/dotfiles/internal/cli"

	_ "yv35.com/dotfiles/internal/logging"
	_ "yv35.com/dotfiles/internal/provider/builtin"
)

// Version information set by ldflags during build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.SetVersionInfo(version, commit, date)
	cli.Execute()
}
