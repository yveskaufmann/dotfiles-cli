package executor

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"yv35.com/dotfiles-cli/internal/theme"
	"yv35.com/dotfiles-cli/internal/tool/git"
)

//go:embed templates/*
var templatesFS embed.FS

// InitializerOptions contains options for the initializer
type InitializerOptions struct {
	TargetDir        string
	GitHubHandle     string
	CreateGitHubRepo bool
	DryRun           bool
}

// Initializer handles scaffolding of new dotfiles repositories
type Initializer struct {
	options InitializerOptions
}

// NewInitializer creates a new initializer
func NewInitializer(options InitializerOptions) *Initializer {
	return &Initializer{
		options: options,
	}
}

// Execute runs the initialization process
func (i *Initializer) Execute() error {
	if i.options.DryRun {
		fmt.Printf("%s[DRY-RUN]%s Running in dry-run mode - no changes will be made\n",
			theme.Colorize(theme.ColorYellow),
			theme.Colorize(theme.ColorReset))
	}

	// Check if already initialized
	if err := i.isAlreadyInitialized(); err != nil {
		return err
	}

	// Create directory structure
	if err := i.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Create files from templates
	if err := i.createFiles(); err != nil {
		return fmt.Errorf("failed to create files: %w", err)
	}

	// Initialize git repository
	if err := i.initializeGitRepository(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Create GitHub repository if requested
	if i.options.CreateGitHubRepo {
		if err := i.createGitHubRepository(); err != nil {
			return fmt.Errorf("failed to create GitHub repository: %w", err)
		}
	}

	fmt.Printf("\n%s✓ Repository initialized successfully%s at: %s\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset),
		i.options.TargetDir)
	fmt.Printf("\n%sNext steps:%s\n", theme.Colorize(theme.ColorCyan), theme.Colorize(theme.ColorReset))
	fmt.Printf("  1. cd %s\n", i.options.TargetDir)
	fmt.Printf("  2. Edit the generated files to customize your setup\n")
	fmt.Printf("  3. git add . && git commit -m 'Initial commit'\n")
	if i.options.CreateGitHubRepo {
		fmt.Printf("  4. git push -u origin main\n")
	}

	return nil
}

func (i *Initializer) isAlreadyInitialized() error {
	if git.IsRepository(i.options.TargetDir) {
		return fmt.Errorf("directory is already a git repository: %s", i.options.TargetDir)
	}

	// Check for key directories
	keyDirs := []string{"bin", "init", "link", "source"}
	for _, dir := range keyDirs {
		path := filepath.Join(i.options.TargetDir, dir)
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("directory already contains dotfiles structure: %s", i.options.TargetDir)
		}
	}

	return nil
}

func (i *Initializer) createDirectoryStructure() error {
	fmt.Printf("%sCreating directory structure...%s\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))

	// Collect unique directories needed from template files
	dirSet := make(map[string]bool)

	// Scan all template files to determine directory structure
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Add subdirectories from templates
			dirSet[entry.Name()] = true
		}
	}

	// Add additional directories that are git-ignored but needed
	dirSet["caches"] = true
	dirSet["logs"] = true

	// Create directories in sorted order for consistent output
	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}

	for _, dir := range dirs {
		path := filepath.Join(i.options.TargetDir, dir)
		if i.options.DryRun {
			fmt.Printf("  [DRY-RUN] mkdir -p %s\n", path)
			continue
		}

		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
		fmt.Printf("  %s✓%s Created %s/\n",
			theme.Colorize(theme.ColorGreen),
			theme.Colorize(theme.ColorReset),
			dir)
	}

	return nil
}

func (i *Initializer) createFiles() error {
	fmt.Printf("%sCreating files from templates...%s\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))

	templateData := struct {
		GitHubHandle  string
		RepositoryURL string
	}{
		GitHubHandle:  i.options.GitHubHandle,
		RepositoryURL: fmt.Sprintf("https://github.com/%s/dotfiles.git", i.options.GitHubHandle),
	}

	// Discover all template files dynamically
	files, err := i.discoverTemplateFiles()
	if err != nil {
		return fmt.Errorf("failed to discover template files: %w", err)
	}

	for targetPath, templatePath := range files {
		if err := i.renderTemplate(templatePath, targetPath, templateData); err != nil {
			return fmt.Errorf("failed to render %s: %w", targetPath, err)
		}

		// Make bin scripts executable
		if strings.HasPrefix(targetPath, "bin/") {
			fullPath := filepath.Join(i.options.TargetDir, targetPath)
			if !i.options.DryRun {
				if err := os.Chmod(fullPath, 0755); err != nil {
					return fmt.Errorf("failed to make %s executable: %w", targetPath, err)
				}
			}
		}

		fmt.Printf("  %s✓%s Created %s\n",
			theme.Colorize(theme.ColorGreen),
			theme.Colorize(theme.ColorReset),
			targetPath)
	}

	return nil
}

// discoverTemplateFiles scans the templates directory and builds a map of target paths to template paths
func (i *Initializer) discoverTemplateFiles() (map[string]string, error) {
	files := make(map[string]string)

	// Read root level templates
	rootEntries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range rootEntries {
		if entry.IsDir() {
			// Read files in subdirectories
			subPath := filepath.Join("templates", entry.Name())
			subEntries, err := templatesFS.ReadDir(subPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", subPath, err)
			}

			for _, subEntry := range subEntries {
				if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".tmpl") {
					// Remove .tmpl extension for target filename
					targetName := strings.TrimSuffix(subEntry.Name(), ".tmpl")
					// Add leading dot for files in link/ directory (these should be dotfiles)
					if entry.Name() == "link" && !strings.HasPrefix(targetName, ".") {
						targetName = "." + targetName
					}
					targetPath := filepath.Join(entry.Name(), targetName)
					templatePath := filepath.Join(subPath, subEntry.Name())
					files[targetPath] = templatePath
				}
			}
		} else if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			// Root level template files
			targetName := strings.TrimSuffix(entry.Name(), ".tmpl")
			templatePath := filepath.Join("templates", entry.Name())
			files[targetName] = templatePath
		}
	}

	return files, nil
}

func (i *Initializer) renderTemplate(templatePath, targetPath string, data interface{}) error {
	// Read template from embedded filesystem
	content, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Parse and execute template
	tmpl, err := template.New(targetPath).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write to target file
	fullPath := filepath.Join(i.options.TargetDir, targetPath)
	if i.options.DryRun {
		fmt.Printf("  [DRY-RUN] write %s\n", fullPath)
		return nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (i *Initializer) initializeGitRepository() error {
	fmt.Printf("%sInitializing git repository...%s\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))

	if i.options.DryRun {
		fmt.Printf("  [DRY-RUN] git init\n")
		fmt.Printf("  [DRY-RUN] git checkout -b main\n")
		return nil
	}

	// Initialize repository
	cmd := exec.Command("git", "init")
	cmd.Dir = i.options.TargetDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git init: %w", err)
	}

	// Create main branch
	cmd = exec.Command("git", "checkout", "-b", "main")
	cmd.Dir = i.options.TargetDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create main branch: %w", err)
	}

	fmt.Printf("  %s✓%s Git repository initialized\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset))
	return nil
}

func (i *Initializer) createGitHubRepository() error {
	fmt.Printf("%sCreating GitHub repository...%s\n",
		theme.Colorize(theme.ColorCyan),
		theme.Colorize(theme.ColorReset))

	// Check if gh CLI is available
	cmd := exec.Command("gh", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Printf("%sWarning:%s GitHub CLI (gh) not found - skipping GitHub repository creation\n",
			theme.Colorize(theme.ColorYellow),
			theme.Colorize(theme.ColorReset))
		fmt.Printf("\nTo create the repository manually:\n")
		fmt.Printf("  1. Create a new repository at https://github.com/new\n")
		fmt.Printf("  2. Run: git remote add origin https://github.com/%s/dotfiles.git\n", i.options.GitHubHandle)
		return nil
	}

	if i.options.DryRun {
		fmt.Printf("  [DRY-RUN] gh repo create %s/dotfiles --public --source=. --remote=origin\n", i.options.GitHubHandle)
		return nil
	}

	// Create repository
	cmd = exec.Command("gh", "repo", "create",
		fmt.Sprintf("%s/dotfiles", i.options.GitHubHandle),
		"--public",
		"--source=.",
		"--remote=origin",
		"--description=Personal dotfiles")
	cmd.Dir = i.options.TargetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create GitHub repository: %w", err)
	}

	fmt.Printf("  %s✓%s GitHub repository created: https://github.com/%s/dotfiles\n",
		theme.Colorize(theme.ColorGreen),
		theme.Colorize(theme.ColorReset),
		i.options.GitHubHandle)
	return nil
}
