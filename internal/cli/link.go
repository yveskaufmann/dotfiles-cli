package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	executor "yv35.com/dotfiles/internal/engine"
	"yv35.com/dotfiles/internal/theme"
	"yv35.com/dotfiles/internal/util/pathutil"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create symlinks from dotfiles to home directory",
	Long: `Link creates symlinks from ~/.dotfiles/link to your home directory.
It handles conflicts and can operate in dry-run mode.`,
	RunE: runLink,
}

var (
	dryRun            bool
	defaultResolution string
)

func init() {
	linkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
	linkCmd.Flags().StringVar(&defaultResolution, "default-resolution", "", "Default conflict resolution: skip, replace, or backup")

	// Register this command with the root command
	RegisterCommands(linkCmd)
}

func runLink(cmd *cobra.Command, args []string) error {
	homeDir := pathutil.MustHomeDir()
	sourceDir := filepath.Join(homeDir, ".dotfiles", "link")
	targetDir := homeDir

	// Parse default resolution
	var resolution executor.LinkConflictResolution
	switch defaultResolution {
	case "skip":
		resolution = executor.ResolutionSkip
	case "replace":
		resolution = executor.ResolutionReplace
	case "backup":
		resolution = executor.ResolutionBackup
	default:
		resolution = executor.ResolutionNone
	}

	linker := executor.NewFileLinker(sourceDir, targetDir, executor.FileLinkerOptions{
		DryRun:                    dryRun,
		DefaultConflictResolution: resolution,
	})

	if dryRun {
		fmt.Println(theme.Colorize(theme.ColorCyan) + "[DRY RUN MODE]" + theme.Colorize(theme.ColorReset))
	}

	return linker.Execute()
}
