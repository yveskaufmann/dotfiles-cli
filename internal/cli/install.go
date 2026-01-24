package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install tools for the specified profile",
	Long:  `Reads the YAML configurations from the init directory and installs all defined tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		initOnly, _ := cmd.Flags().GetBool("init-only")
		if err := runInitScripts(cmd, initOnly); err != nil {
			fmt.Printf("❌ Installation failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RegisterCommands(installCmd)
}
