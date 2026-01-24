package setuptools

import (
	"fmt"

	_os "yv35.com/dotfiles/internal/os"
	"yv35.com/dotfiles/internal/tools"
)

type Executor struct {
	Groups []Group
}

func NewExecutor(groups []Group) *Executor {
	return &Executor{Groups: groups}
}

func (e *Executor) Execute() error {
	for _, group := range e.Groups {
		if group.Systems != "" && !_os.Is(_os.OSType(group.Systems)) {
			// fmt.Printf("⏭️  Skipping group %s (system mismatch: %s)\n", group.Name, group.Systems)
			continue
		}

		fmt.Printf("🚀 Processing group: %s\n", group.Name)
		if err := e.executeGroup(group); err != nil {
			return fmt.Errorf("failed to execute group %s: %w", group.Name, err)
		}
	}
	return nil
}

func (e *Executor) executeGroup(group Group) error {
	// 1. Apt packages
	if len(group.Apt) > 0 {
		if _os.IsLinux() {
			if err := tools.InstallAptPackages(group.Apt); err != nil {
				return err
			}
		}
	}

	// 2. Snap packages
	if len(group.Snap) > 0 {
		if _os.IsUbuntu() {
			var snaps []string
			for _, s := range group.Snap {
				// For now just name, ignoring classic for a moment or passing it
				snaps = append(snaps, s.Name)
			}
			if err := tools.InstallSnapPackages(snaps); err != nil {
				return err
			}
		}
	}

	// 3. Pipx packages
	if len(group.Pipx) > 0 {
		for _, p := range group.Pipx {
			if err := tools.InstallPipx(p.Name); err != nil {
				return err
			}
		}
	}

	// 4. NPM packages
	if len(group.NPM) > 0 {
		for _, n := range group.NPM {
			if err := tools.InstallNPM(n.Name); err != nil {
				return err
			}
		}
	}

	// 5. Github Releases
	if len(group.Github) > 0 {
		for _, g := range group.Github {
			if err := tools.InstallGithubRelease(g.Name, g.Repo, g.Version, g.AssetPattern, g.BinaryPath, g.InstallPath, g.Binaries); err != nil {
				return err
			}
		}
	}

	// 6. Binary downloads
	if len(group.Binary) > 0 {
		for _, b := range group.Binary {
			if err := tools.InstallBinary(b.Name, b.URL, b.Version, b.BinaryPath, b.InstallPath, b.Binaries); err != nil {
				return err
			}
		}
	}

	// 7. PPA
	if len(group.PPA) > 0 {
		if _os.IsLinux() {
			for _, p := range group.PPA {
				if err := tools.InstallPPA(p.Name, p.Key, p.KeyServer, p.KeyID, p.URI, p.Suites, p.Components, p.Pkgs); err != nil {
					return err
				}
			}
		}
	}

	// 7. Custom scripts
	for _, c := range group.Custom {
		if err := e.executeCustom(c); err != nil {
			return err
		}
	}

	// 8. Scripts
	for _, s := range group.Script {
		if err := tools.RunScript(s.Script); err != nil {
			return fmt.Errorf("failed to run script %s: %w", s.Name, err)
		}
	}

	return nil
}

func (e *Executor) executeCustom(c CustomSpec) error {
	if c.InstallCheck != "" {
		if err := tools.RunShell(c.InstallCheck); err == nil {
			fmt.Printf("✅ %s is already installed (check passed)\n", c.Name)
			return nil
		}
	}

	fmt.Printf("⚙️  Installing %s...\n", c.Name)
	if err := tools.RunShell(c.Install); err != nil {
		return fmt.Errorf("failed to install %s: %w", c.Name, err)
	}
	fmt.Printf("✅ %s installed successfully\n", c.Name)
	return nil
}
