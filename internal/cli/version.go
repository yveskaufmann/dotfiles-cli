package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionInfo = struct {
		version string
		commit  string
		date    string
	}{
		version: "dev",
		commit:  "none",
		date:    "unknown",
	}
)

// SetVersionInfo sets the version information from main package
func SetVersionInfo(version, commit, date string) {
	versionInfo.version = version
	versionInfo.commit = commit
	versionInfo.date = date
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version, commit hash, and build date of the dotfiles CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s\n", versionInfo.version)
		fmt.Printf("commit: %s\n", versionInfo.commit)
		fmt.Printf("built:  %s\n", versionInfo.date)
	},
}

func init() {
	RegisterCommands(versionCmd)
}
