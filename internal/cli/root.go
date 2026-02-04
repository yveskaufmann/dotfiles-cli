/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"yv35.com/dotfiles/internal/tool/git"
	"yv35.com/dotfiles/internal/util/fsutil"
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

func init() {
	rootCmd.Flags().BoolP("pull-only", "p", false, "Only pull dotfiles changes from the remote repository.")
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

	if false {
		// 2. Run Init scripts
		if err := runInitScripts(cmd, initOnly); err != nil {
			fmt.Println("failed to run init scripts:", err)
			os.Exit(1)
		}
	}

	// 3. Link system files
	if !initOnly && false {
		if err := linkSystemFiles(); err != nil {
			fmt.Println("failed to link system files:", err)
			os.Exit(1)
		}
	}

	if false {
		// 4. Copy Dotfiles into home directory
		if err := copyDotfiles(pullOnly); err != nil {
			fmt.Println("failed to copy dotfiles:", err)
			os.Exit(1)
		}
	}

	fmt.Println("✅ Bootstrap completed successfully.")
}

func linkSystemFiles() error {
	return nil
}

func copyDotfiles(pullOnly bool) error {
	panic("unimplemented")
}
