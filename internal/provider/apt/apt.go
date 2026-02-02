package apt

import (
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"strings"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/types"
	"yv35.com/dotfiles/internal/util/osutil"
	"yv35.com/dotfiles/internal/util/sh"
)

type Provider struct {
	aptUpdated bool
}

func NewProvider() *Provider {
	return &Provider{
		aptUpdated: false,
	}
}

func (p *Provider) ID() string {
	return "apt"
}

func (p *Provider) Priority() int {
	return 10 // System package manager - high priority
}

func (p *Provider) Setup() error {
	// Skip if not Linux
	if !osutil.IsLinux() {
		return nil
	}

	// Ensure apt update is run once
	return p.ensureAptUpdate(false)
}

func (p *Provider) Install(group config.DependencyGroup, onComplete types.OnTaskComplete) error {
	// Skip if not Linux
	if !osutil.IsLinux() {
		return nil
	}

	var allPackagesToInstall []string
	ppaSourcesAdded := false

	// Phase 1: Add all PPA sources (no update yet)
	for _, ppa := range group.PPA {

		slog.Info("Processing PPA source", "name", ppa.Name, "uri", ppa.URI)

		added, packages, err := p.addPPASource(ppa, onComplete)
		if err != nil {
			return err
		}
		if added {
			ppaSourcesAdded = true
		}
		// Collect packages from this PPA
		for _, pkg := range packages {
			if !p.isAptInstalled(pkg) {
				allPackagesToInstall = append(allPackagesToInstall, pkg)
			}
		}
	}

	// Phase 2: Single apt-update if any PPAs were added
	if ppaSourcesAdded {
		if err := p.ensureAptUpdate(true); err != nil {
			return err
		}
	}

	// Phase 3: Process regular Apt packages
	for _, pkg := range group.Apt {
		if p.isAptInstalled(pkg) {
			onComplete(types.TaskResult{
				Name:   pkg,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already installed"),
			})
			continue
		}

		if p.isURL(pkg) {
			if err := p.installDebFromURL(pkg, onComplete); err != nil {
				return err
			}
		} else {
			allPackagesToInstall = append(allPackagesToInstall, pkg)
		}
	}

	// Phase 4: Batch install all packages at once
	if len(allPackagesToInstall) > 0 {
		if err := p.batchInstallWithVerification(allPackagesToInstall, onComplete); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) isAptInstalled(pkg string) bool {

	if p.isURL(pkg) {
		pkg = p.getPkgNameFromURL(pkg)
		if pkg == "" {
			// When we can't determine the package name from the URL, assume not installed
			return false
		}
	}

	cmd := fmt.Sprintf("dpkg-query -W --showformat='${Status}\\n' %s 2>/dev/null | grep 'install ok installed'", pkg)
	return sh.RunShell(cmd) == nil
}

func (p *Provider) getPkgNameFromURL(pkg string) string {
	u, err := url.Parse(pkg)
	if err != nil {
		return ""
	}

	fileName := path.Base(u.Path)
	if fileName == "." || fileName == "/" {
		return ""
	}

	components := strings.FieldsFunc(fileName, func(r rune) bool {
		return r == '_' || r == '-'
	})

	if len(components) == 0 {
		return ""
	}

	return components[0]
}
