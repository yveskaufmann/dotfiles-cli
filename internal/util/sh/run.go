package sh

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func RunShell(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	c := exec.CommandContext(ctx, "bash", "-lc", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func RunShellOutput(cmd string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	c := exec.CommandContext(ctx, "bash", "-lc", cmd)
	b, err := c.Output()
	return string(b), err
}

func RunScript(path string) error {
	// If path is relative, try to resolve it relative to $DOTFILES
	if !filepath.IsAbs(path) {
		dotfiles := os.Getenv("DOTFILES")
		if dotfiles != "" {
			fullPath := filepath.Join(dotfiles, path)
			if _, err := os.Stat(fullPath); err == nil {
				path = fullPath
			}
		}
	}
	return RunShell(path)
}
