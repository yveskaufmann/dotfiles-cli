package main

import (
	"fmt"
	"os"
	"path"

	executor "yv35.com/dotfiles/internal/engine"
)

func main() {

	userHomeDir := MustHomeDir()
	dotfilesLinkDirectory := path.Join(userHomeDir, ".dotfiles", "link")

	fmt.Printf("Dotfile path: %s\n", dotfilesLinkDirectory)

	linker := executor.NewFileLinker(dotfilesLinkDirectory, "/tmp/file-traverse1419341182")
	if err := linker.Execute(); err != nil {
		fmt.Printf("❌ Failed to link dotfiles to home directory: %v\n", err)
	}

	fmt.Println("✅  Dotfiles linked successfully.")
}

func MustHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get user home directory: %w", err))
	}
	return homeDir
}
