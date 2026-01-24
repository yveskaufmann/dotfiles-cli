package setuptools

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	InitDir string
	Profile string
}

func NewLoader(initDir string, profile string) *Loader {
	return &Loader{
		InitDir: initDir,
		Profile: profile,
	}
}

func (l *Loader) Load() ([]Group, error) {
	files, err := filepath.Glob(filepath.Join(l.InitDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list yaml files in %s: %w", l.InitDir, err)
	}

	sort.Strings(files)

	var allGroups []Group
	for _, file := range files {
		groups, err := l.loadFile(file)
		if err != nil {
			return nil, err
		}
		allGroups = append(allGroups, groups...)
	}

	return l.filterGroups(allGroups), nil
}

func (l *Loader) loadFile(path string) ([]Group, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml in %s: %w", path, err)
	}

	return config.Groups, nil
}

func (l *Loader) filterGroups(groups []Group) []Group {
	var filtered []Group
	for _, g := range groups {
		if g.Profile == "" || g.Profile == l.Profile || g.Profile == "default" {
			filtered = append(filtered, g)
		}
	}
	return filtered
}
