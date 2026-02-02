package main

import (
	"yv35.com/dotfiles/internal/cli"

	_ "yv35.com/dotfiles/internal/logging"
	_ "yv35.com/dotfiles/internal/provider/builtin"
)

func main() {
	cli.Execute()
}
