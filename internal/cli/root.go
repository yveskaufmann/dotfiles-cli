/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"yv35.com/dotfiles/internal/fsutil"
	setuptools "yv35.com/dotfiles/internal/init/setup-tools"
	"yv35.com/dotfiles/internal/tools/git"
)

var rootCmd = &cobra.Command{

	Use:   "dotfiles",
	Short: "Bootstrap dotfiles, install packages, and create symlinks",
	Long: `DOTFILES - Yves Kaufmann - https://github.com/yveskaufmann/dotfiles
	
This CLI tool bootstraps my dotfiles by symlinking files to the home directory
and installing/updating required packages.`,

	Example: ` 
  # Bootstrap the .dotfiles and install dependencies
  $ dotfiles 
  
  # Install / Update dependencies
  $ dotfiles --init-only`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommands are provided
		Bootstrap(cmd)
	},
}

func Init() {
	rootCmd.PersistentFlags().BoolP("init-only", "i", false, "Executes only the init scripts steps, skipping symlink creation.")
	rootCmd.Flags().BoolP("pull-only", "p", false, "Only pull dotfiles changes from the remote repository.")
	rootCmd.PersistentFlags().StringP("profile", "P", "default", "The setup profile to use (e.g., 'default', 'work').")

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// RegisterCommands accepts a list of commands to add as subcommands of rootCmd
func RegisterCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		rootCmd.AddCommand(cmd)
	}
}

func Bootstrap(cmd *cobra.Command) {
	initOnly, _ := cmd.Flags().GetBool("init-only")
	pullOnly, _ := cmd.Flags().GetBool("pull-only")

	if err := git.Ensure(); err != nil {
		fmt.Println("failed to bootstrap: git is missing, and could not be installed automatically", err)
		os.Exit(1)
	}

	if err := fsutil.EnsureDirectory(CACHE_PATH); err != nil {
		fmt.Println("failed to create caches directory:", err)
		os.Exit(1)
	}

	// 2. Run Init scripts
	if err := runInitScripts(cmd, initOnly); err != nil {
		fmt.Println("failed to run init scripts:", err)
		os.Exit(1)
	}

	// 3. Link system files
	if !initOnly {
		if err := linkSystemFiles(); err != nil {
			fmt.Println("failed to link system files:", err)
			os.Exit(1)
		}
	}

	// 4. Copy Dotfiles into home directory
	if err := copyDotfiles(pullOnly); err != nil {
		fmt.Println("failed to copy dotfiles:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Bootstrap completed successfully.")
}

func runInitScripts(cmd *cobra.Command, initOnly bool) error {
	profile, _ := cmd.Flags().GetString("profile")

	fmt.Printf("📦 Initializing system (Profile: %s)...\n", profile)

	loader := setuptools.NewLoader(INIT_SCRIPTS_PATH, profile)
	groups, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load tool configurations: %w", err)
	}

	executor := setuptools.NewExecutor(groups)
	if err := executor.Execute(); err != nil {
		return fmt.Errorf("bootstrap failed during tool installation: %w", err)
	}

	fmt.Println("✅ System tools initialization completed.")
	fmt.Println("👉 Some tools might require a shell restart or 'source ~/.zshrc' to be available on PATH.")
	return nil
}

func linkSystemFiles() error {
	panic("unimplemented")
}

func copyDotfiles(pullOnly bool) error {
	panic("unimplemented")
}
