package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "Bootstrap dotfiles, install packages, and create symlinks",
	Long: `DOTFILES CLI - https://github.com/yveskaufmann/dotfiles-cli
	
This CLI tool bootstraps dotfiles by symlinking files to the home directory
and installing/updating required packages.`,

	Example: `  # Bootstrap your dotfiles (clone/pull repository, install tools, create symlinks)
  $ dotfiles bootstrap
  
  # Install tools for the default profile
  $ dotfiles install
  
  # Create symlinks from dotfiles to home directory
  $ dotfiles link
  
  # Bootstrap with custom repository
  $ dotfiles bootstrap --repository git@github.com:user/dotfiles.git
  
  # Bootstrap for work profile
  $ dotfiles bootstrap --profile work`,
}

func init() {
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
