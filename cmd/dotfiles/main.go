package main

import (
	"yv35.com/dotfiles/internal/cli"
	"yv35.com/dotfiles/internal/logging"
)

func main() {
	logging.Init()

	cli.Execute()
}
