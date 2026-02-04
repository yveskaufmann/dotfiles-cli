package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"yv35.com/dotfiles/internal/config"
	executor "yv35.com/dotfiles/internal/engine"
	"yv35.com/dotfiles/internal/theme"
	"yv35.com/dotfiles/internal/tool/git"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap dotfiles by cloning repository, installing tools, and creating symlinks",
	Long: `Bootstrap sets up your dotfiles environment from scratch or updates an existing setup.

This command will:
1. Clone your dotfiles repository (or pull latest changes if it exists)
2. Install tools and packages defined in your init scripts
3. Create symlinks from your dotfiles to your home directory

By default, all steps are executed. Use flags to run only specific steps or skip certain steps.`,
	Example: `  # Full bootstrap (clone/pull, install, link)
  $ dotfiles bootstrap

  # Bootstrap with custom repository
  $ dotfiles bootstrap --repository git@github.com:user/dotfiles.git

  # Only install tools (skip linking)
  $ dotfiles bootstrap --install

  # Only create symlinks (skip install)
  $ dotfiles bootstrap --link

  # Bootstrap without pulling latest changes
  $ dotfiles bootstrap --no-pull

  # Bootstrap with specific profile and dry-run for linking
  $ dotfiles bootstrap --profile work --dry-run`,
	RunE: runBootstrap,
}

var (
	bootstrapRepository     string
	bootstrapNoPull         bool
	bootstrapNoInstall      bool
	bootstrapNoLink         bool
	bootstrapInstallOnly    bool
	bootstrapLinkOnly       bool
	bootstrapProfile        string
	bootstrapProviders      []string
	bootstrapLinkDryRun     bool
	bootstrapLinkResolution string
)

func init() {
	// Repository options
	bootstrapCmd.Flags().StringVar(&bootstrapRepository, "repository", "", "Git repository URL (overrides config)")
	bootstrapCmd.Flags().BoolVar(&bootstrapNoPull, "no-pull", false, "Skip pulling latest changes if repository exists")

	// Step control flags
	bootstrapCmd.Flags().BoolVar(&bootstrapInstallOnly, "install", false, "Only run install step (skip linking)")
	bootstrapCmd.Flags().BoolVar(&bootstrapLinkOnly, "link", false, "Only run link step (skip install)")
	bootstrapCmd.Flags().BoolVar(&bootstrapNoInstall, "no-install", false, "Skip install step")
	bootstrapCmd.Flags().BoolVar(&bootstrapNoLink, "no-link", false, "Skip link step")

	// Install options
	bootstrapCmd.Flags().StringVarP(&bootstrapProfile, "profile", "P", "default", "The setup profile to use (e.g., 'default', 'work')")
	bootstrapCmd.Flags().StringSliceVar(&bootstrapProviders, "providers", []string{}, "Comma-separated list of providers to enable (e.g., 'nvm,apt')")

	// Link options
	bootstrapCmd.Flags().BoolVar(&bootstrapLinkDryRun, "dry-run", false, "Show what would be done without making changes (for linking)")
	bootstrapCmd.Flags().StringVar(&bootstrapLinkResolution, "default-resolution", "", "Default conflict resolution: skip, replace, or backup")

	RegisterCommands(bootstrapCmd)
}

func runBootstrap(cmd *cobra.Command, args []string) error {
	// Determine which steps to run
	runInstall, runLink := determineSteps()
	runPull := !bootstrapNoPull

	// Load or prompt for repository URL
	repositoryURL, err := getRepositoryURL()
	if err != nil {
		return err
	}

	// Create config save callback
	saveConfigFn := func() error {
		cfg := &config.BootstrapConfig{}
		cfg.Dotfiles.Repository.URL = repositoryURL
		cfg.Dotfiles.Repository.Type = detectRepositoryType(repositoryURL)
		return config.SaveBootstrapConfig(cfg)
	}

	// Create bootstrapper
	opts := executor.BootstrapOptions{
		RunPull:        runPull,
		RunInstall:     runInstall,
		RunLink:        runLink,
		Profile:        bootstrapProfile,
		Providers:      bootstrapProviders,
		LinkDryRun:     bootstrapLinkDryRun,
		LinkResolution: bootstrapLinkResolution,
		SaveConfig:     saveConfigFn,
	}

	bootstrapper := executor.NewBootstrapper(repositoryURL, opts)

	// Execute bootstrap with retry on clone failure
	err = bootstrapper.Execute()
	if err != nil && strings.Contains(err.Error(), "failed to clone repository") {
		// Offer to change repository URL
		fmt.Printf("\n%s⚠️  Clone failed: %v%s\n\n",
			theme.Colorize(theme.ColorRed),
			err,
			theme.Colorize(theme.ColorReset))

		newURL, promptErr := promptRepositoryURL(repositoryURL, true)
		if promptErr != nil {
			return err // Return original error
		}

		if newURL != "" && newURL != repositoryURL {
			// Save new URL to config
			cfg := &config.BootstrapConfig{}
			cfg.Dotfiles.Repository.URL = newURL
			cfg.Dotfiles.Repository.Type = detectRepositoryType(newURL)
			if saveErr := config.SaveBootstrapConfig(cfg); saveErr != nil {
				fmt.Printf("%s⚠️  Warning: Failed to save config: %v%s\n",
					theme.Colorize(theme.ColorYellow),
					saveErr,
					theme.Colorize(theme.ColorReset))
			}

			// Retry with new URL - need to update saveConfigFn with new URL
			newSaveConfigFn := func() error {
				cfg := &config.BootstrapConfig{}
				cfg.Dotfiles.Repository.URL = newURL
				cfg.Dotfiles.Repository.Type = detectRepositoryType(newURL)
				return config.SaveBootstrapConfig(cfg)
			}
			opts.SaveConfig = newSaveConfigFn
			bootstrapper = executor.NewBootstrapper(newURL, opts)
			return bootstrapper.Execute()
		}

		return err
	}

	return err
}

// determineSteps determines which bootstrap steps should run based on flags
func determineSteps() (runInstall, runLink bool) {
	// If --install or --link specified, only run those steps
	if bootstrapInstallOnly || bootstrapLinkOnly {
		return bootstrapInstallOnly, bootstrapLinkOnly
	}

	// Otherwise, run all steps except those explicitly skipped
	runInstall = !bootstrapNoInstall
	runLink = !bootstrapNoLink

	return runInstall, runLink
}

// getRepositoryURL determines the repository URL from flag, config, or prompt
func getRepositoryURL() (string, error) {
	// 1. Check flag
	if bootstrapRepository != "" {
		return bootstrapRepository, nil
	}

	// 2. Check if repository already exists and get URL from git remote
	if git.IsRepository(DOTFILES_PATH) {
		if remoteURL := git.GetRemoteURL(DOTFILES_PATH); remoteURL != "" {
			return remoteURL, nil
		}
	}

	// 3. Check config file
	cfg, err := config.LoadBootstrapConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	if cfg != nil && cfg.Dotfiles.Repository.URL != "" {
		return cfg.Dotfiles.Repository.URL, nil
	}

	// 4. Prompt user
	return promptRepositoryURL("", false)
}

// promptRepositoryURL prompts the user for a repository URL
func promptRepositoryURL(currentURL string, isRetry bool) (string, error) {
	if isRetry {
		fmt.Printf("%s⚠️  The repository could not be cloned.%s\n",
			theme.Colorize(theme.ColorYellow),
			theme.Colorize(theme.ColorReset))
		fmt.Printf("Would you like to try a different repository URL?\n\n")
	} else {
		fmt.Printf("%sNo repository configured.%s\n",
			theme.Colorize(theme.ColorYellow),
			theme.Colorize(theme.ColorReset))
		fmt.Printf("Please enter your dotfiles repository URL:\n\n")
	}

	if currentURL != "" {
		fmt.Printf("Current: %s\n", currentURL)
	}
	fmt.Printf("Example: git@github.com:username/dotfiles.git\n")
	fmt.Printf("         https://github.com/username/dotfiles.git\n\n")
	fmt.Printf("Repository URL: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	url := strings.TrimSpace(input)
	if url == "" {
		if currentURL != "" {
			// Keep current URL
			return currentURL, nil
		}
		return "", fmt.Errorf("repository URL is required")
	}

	return url, nil
}

// detectRepositoryType detects if the repository URL is ssh or https
func detectRepositoryType(url string) string {
	if strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://") {
		return "ssh"
	}
	return "https"
}
