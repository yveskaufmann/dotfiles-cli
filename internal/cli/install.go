package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"yv35.com/dotfiles-cli/internal/config"
	executor "yv35.com/dotfiles-cli/internal/engine"
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

func runInitScripts(cmd *cobra.Command, initOnly bool) error {
	profile, _ := cmd.Flags().GetString("profile")
	enabledProviders, _ := cmd.Flags().GetStringSlice("providers")

	fmt.Printf("📦 Initializing system (Profile: %s)...\n", profile)

	loader := config.NewLoader(INIT_SCRIPTS_PATH, profile)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load tool configurations: %w", err)
	}

	executor := executor.NewToolInstaller(cfg)
	executor.SetEnabledProviders(enabledProviders)

	if err := executor.Execute(); err != nil {
		return fmt.Errorf("bootstrap failed during tool installation: %w", err)
	}

	fmt.Println("✅ System tools initialization completed.")
	fmt.Println("👉 Some tools might require a shell restart or 'source ~/.zshrc' to be available on PATH.")
	return nil
}

func init() {
	installCmd.Flags().BoolP("init-only", "i", false, "Executes only the init scripts steps, skipping symlink creation.")
	installCmd.Flags().StringP("profile", "P", "default", "The setup profile to use (e.g., 'default', 'work').")
	installCmd.Flags().StringSliceP("providers", "", []string{}, "Comma-separated list of providers to enable (e.g., 'nvm,apt'). If empty, all providers are enabled.")
	RegisterCommands(installCmd)
}
