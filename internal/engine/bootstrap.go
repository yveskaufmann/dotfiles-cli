package executor

import (
	"fmt"
	"path/filepath"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/theme"
	"yv35.com/dotfiles/internal/tool/git"
	"yv35.com/dotfiles/internal/util/fsutil"
	"yv35.com/dotfiles/internal/util/pathutil"
)

// BootstrapOptions contains configuration options for the Bootstrapper
type BootstrapOptions struct {
	RunPull        bool
	RunInstall     bool
	RunLink        bool
	Profile        string
	Providers      []string
	LinkDryRun     bool
	LinkResolution string
	DotfilesPath   string
}

// Bootstrapper orchestrates the bootstrap process: clone/pull dotfiles, install tools, create symlinks
type Bootstrapper struct {
	repositoryURL  string
	dotfilesPath   string
	runPull        bool
	runInstall     bool
	runLink        bool
	profile        string
	providers      []string
	linkDryRun     bool
	linkResolution LinkConflictResolution
}

// NewBootstrapper creates a new Bootstrapper instance
func NewBootstrapper(repositoryURL string, opts BootstrapOptions) *Bootstrapper {
	dotfilesPath := opts.DotfilesPath
	if dotfilesPath == "" {
		homeDir := pathutil.MustHomeDir()
		dotfilesPath = filepath.Join(homeDir, ".dotfiles")
	}

	// Parse link resolution
	var resolution LinkConflictResolution
	switch opts.LinkResolution {
	case "skip":
		resolution = ResolutionSkip
	case "replace":
		resolution = ResolutionReplace
	case "backup":
		resolution = ResolutionBackup
	default:
		resolution = ResolutionNone
	}

	return &Bootstrapper{
		repositoryURL:  repositoryURL,
		dotfilesPath:   dotfilesPath,
		runPull:        opts.RunPull,
		runInstall:     opts.RunInstall,
		runLink:        opts.RunLink,
		profile:        opts.Profile,
		providers:      opts.Providers,
		linkDryRun:     opts.LinkDryRun,
		linkResolution: resolution,
	}
}

// Execute runs the bootstrap process
func (b *Bootstrapper) Execute() error {
	fmt.Printf("%s🚀 Starting bootstrap process...%s\n\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))

	// 1. Ensure git is installed
	fmt.Printf("%s📋 Step 1/4: Ensuring git is installed...%s\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))
	if err := git.Ensure(); err != nil {
		return fmt.Errorf("failed to ensure git is installed: %w", err)
	}
	fmt.Printf("%s✅ Git is available%s\n\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset))

	// 2. Ensure repository (clone or pull)
	fmt.Printf("%s📋 Step 2/4: Setting up dotfiles repository...%s\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))
	if err := b.EnsureRepository(); err != nil {
		return fmt.Errorf("failed to setup repository: %w", err)
	}
	fmt.Printf("%s✅ Repository ready at %s%s\n\n",
		theme.Colorize(theme.ColorGreen),
		pathutil.MinimizePath(b.dotfilesPath),
		theme.Colorize(theme.ColorReset))

	// 3. Run install if requested
	if b.runInstall {
		fmt.Printf("%s📋 Step 3/4: Installing tools...%s\n",
			theme.Colorize(theme.ColorCyan),
			theme.Colorize(theme.ColorReset))
		if err := b.RunInstall(); err != nil {
			return fmt.Errorf("failed to install tools: %w", err)
		}
		fmt.Printf("%s✅ Tools installed%s\n\n",
			theme.Colorize(theme.ColorGreen),
			theme.Colorize(theme.ColorReset))
	} else {
		fmt.Printf("%s⏭️  Step 3/4: Skipping tool installation%s\n\n",
			theme.Colorize(theme.ColorYellow),
			theme.Colorize(theme.ColorReset))
	}

	// 4. Run link if requested
	if b.runLink {
		fmt.Printf("%s📋 Step 4/4: Creating symlinks...%s\n",
			theme.Colorize(theme.ColorCyan),
			theme.Colorize(theme.ColorReset))
		if err := b.RunLink(); err != nil {
			return fmt.Errorf("failed to create symlinks: %w", err)
		}
		fmt.Printf("%s✅ Symlinks created%s\n\n",
			theme.Colorize(theme.ColorGreen),
			theme.Colorize(theme.ColorReset))
	} else {
		fmt.Printf("%s⏭️  Step 4/4: Skipping symlink creation%s\n\n",
			theme.Colorize(theme.ColorYellow),
			theme.Colorize(theme.ColorReset))
	}

	// Print success message
	fmt.Printf("%s╔════════════════════════════════════════════════════════╗%s\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset))
	fmt.Printf("%s║  ✅ Bootstrap completed successfully!                  ║%s\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset))
	fmt.Printf("%s╚════════════════════════════════════════════════════════╝%s\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset))

	return nil
}

// EnsureRepository ensures the dotfiles repository exists
// Clones if missing, pulls if exists and runPull is true
func (b *Bootstrapper) EnsureRepository() error {
	if git.IsRepository(b.dotfilesPath) {
		fmt.Printf("Repository already exists at %s\n", pathutil.MinimizePath(b.dotfilesPath))

		if b.runPull {
			fmt.Printf("Pulling latest changes...\n")
			if err := git.Pull(b.dotfilesPath); err != nil {
				return fmt.Errorf("failed to pull repository: %w", err)
			}
			fmt.Printf("Successfully updated repository\n")
		} else {
			fmt.Printf("Skipping pull (--no-pull specified)\n")
		}
		return nil
	}

	// Repository doesn't exist, clone it
	fmt.Printf("Cloning repository from %s...\n", b.repositoryURL)
	if err := git.Clone(b.repositoryURL, b.dotfilesPath); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	fmt.Printf("Successfully cloned repository\n")

	// Ensure cache directory exists
	cacheDir := filepath.Join(b.dotfilesPath, ".caches")
	if err := fsutil.EnsureDirectory(cacheDir); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	return nil
}

// RunInstall executes the tool installation process
func (b *Bootstrapper) RunInstall() error {
	initScriptsPath := filepath.Join(b.dotfilesPath, "init")

	loader := config.NewLoader(initScriptsPath, b.profile)
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load tool configurations: %w", err)
	}

	installer := NewToolInstaller(cfg)
	installer.SetEnabledProviders(b.providers)

	if err := installer.Execute(); err != nil {
		return fmt.Errorf("tool installation failed: %w", err)
	}

	return nil
}

// RunLink executes the symlink creation process
func (b *Bootstrapper) RunLink() error {
	sourceDir := filepath.Join(b.dotfilesPath, "link")
	targetDir := pathutil.MustHomeDir()

	linker := NewFileLinker(sourceDir, targetDir, FileLinkerOptions{
		DryRun:                    b.linkDryRun,
		DefaultConflictResolution: b.linkResolution,
	})

	if b.linkDryRun {
		fmt.Printf("%s[DRY RUN MODE]%s\n",
			theme.Colorize(theme.ColorCyan),
			theme.Colorize(theme.ColorReset))
	}

	return linker.Execute()
}
