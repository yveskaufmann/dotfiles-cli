package jetbrains

import (
	"fmt"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
)

// installCandidate bundles a spec together with the already-resolved release
// info so that the worker does not have to fetch it a second time.
type installCandidate struct {
	Spec    config.JetbrainsSpec
	Release *ReleaseInfo
	Dir     string
}

// checkNeedsInstall resolves the desired release and compares it with the
// version that is currently installed at /opt/<ide-dir>.
// Returns (candidate, true, nil) when an install/update is required.
// Returns (zero, false, nil) when the installed version already matches.
func (p *Provider) checkNeedsInstall(spec config.JetbrainsSpec) (installCandidate, bool, error) {
	code, err := resolveCode(spec.IDE)
	if err != nil {
		return installCandidate{}, false, err
	}

	dir, err := installDir(code)
	if err != nil {
		return installCandidate{}, false, err
	}

	release, err := FetchRelease(code, spec.Version)
	if err != nil {
		return installCandidate{}, false, err
	}

	installedVersion, err := getInstalledVersion(dir)
	if err != nil {
		return installCandidate{}, false, fmt.Errorf("could not read installed version of %s: %w", spec.IDE, err)
	}

	if installedVersion == release.Version {
		return installCandidate{}, false, nil
	}

	return installCandidate{Spec: spec, Release: release, Dir: dir}, true, nil
}

// filterInstalled calls onComplete(StatusSkipped) for IDEs that are already at
// the desired version, and returns the remaining specs that need installation.
func (p *Provider) filterInstalled(specs []config.JetbrainsSpec, onComplete types.OnTaskComplete) []installCandidate {
	var toInstall []installCandidate

	for _, spec := range specs {
		displayName := spec.IDE
		if spec.Name != "" {
			displayName = spec.Name
		}

		candidate, needsInstall, err := p.checkNeedsInstall(spec)
		if err != nil {
			// Surface the error as a failed task so the engine can report it.
			onComplete(types.TaskResult{
				Name:   displayName,
				Status: types.StatusFailed,
				Error:  err,
			})
			continue
		}

		if !needsInstall {
			onComplete(types.TaskResult{
				Name:   displayName,
				Status: types.StatusSkipped,
				Error:  fmt.Errorf("already up-to-date"),
			})
			continue
		}

		toInstall = append(toInstall, candidate)
	}

	return toInstall
}
