package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	executor "yv35.com/dotfiles-cli/internal/engine"
)

var (
	initDryRun bool
	initGitHub bool
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new dotfiles repository",
	Long: `Initialize a new dotfiles repository with a standard structure.

This command scaffolds a new dotfiles repository with the following structure:
  - bin/       Custom scripts and utilities
  - init/      Tool installation configurations
  - link/      Files to be symlinked to home directory
  - source/    Shell scripts to be sourced
  - caches/    Cache directory (git-ignored)
  - logs/      Log directory (git-ignored)

The repository will be initialized as a git repository with a main branch.`,
	Example: `  # Initialize in current directory
  dotfiles init

  # Initialize in specific directory
  dotfiles init ~/my-dotfiles

  # Preview changes without creating files
  dotfiles init --dry-run

  # Initialize and create GitHub repository
  dotfiles init --github`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Preview changes without creating files")
	initCmd.Flags().BoolVar(&initGitHub, "github", false, "Create a GitHub repository (requires gh CLI)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Create target directory if it doesn't exist
	if !initDryRun {
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	// Prompt for GitHub handle if creating GitHub repo
	githubHandle := ""
	if initGitHub || !initDryRun {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("GitHub username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read GitHub handle: %w", err)
		}
		githubHandle = strings.TrimSpace(input)
		if githubHandle == "" {
			return fmt.Errorf("GitHub username is required")
		}
	}

	// Create initializer
	initializer := executor.NewInitializer(executor.InitializerOptions{
		TargetDir:        absPath,
		GitHubHandle:     githubHandle,
		CreateGitHubRepo: initGitHub,
		DryRun:           initDryRun,
	})

	// Execute initialization
	return initializer.Execute()
}
