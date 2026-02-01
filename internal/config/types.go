package config

import "gopkg.in/yaml.v3"

type Config struct {
	Schema string            `yaml:"$schema,omitempty"`
	Groups []DependencyGroup `yaml:"groups"`
}

type DependencyGroup struct {
	Name        string        `yaml:"name"`
	Profile     string        `yaml:"profile,omitempty"`
	Systems     string        `yaml:"systems,omitempty"`
	Apt         []string      `yaml:"apt,omitempty"`
	PPA         []PPASpec     `yaml:"ppa,omitempty"`
	Snap        []SnapSpec    `yaml:"snap,omitempty"`
	Pipx        []PipxSpec    `yaml:"pipx,omitempty"`
	NPM         []NPMSpec     `yaml:"npm,omitempty"`
	Brew        []BrewSpec    `yaml:"brew,omitempty"`
	BrewTapSpec []BrewTapSpec `yaml:"brew_taps,omitempty"`
	Github      []GithubSpec  `yaml:"github_release,omitempty"`
	Binary      []BinarySpec  `yaml:"binary,omitempty"`
	Script      []ScriptSpec  `yaml:"script,omitempty"`
	Custom      []CustomSpec  `yaml:"custom,omitempty"`
}

type PPASpec struct {
	Name       string   `yaml:"name"`
	Key        string   `yaml:"key,omitempty"`
	URI        string   `yaml:"uri,omitempty"`
	Suites     string   `yaml:"suites,omitempty"`
	Components string   `yaml:"components,omitempty"`
	KeyServer  string   `yaml:"key_server,omitempty"`
	KeyID      string   `yaml:"key_id,omitempty"`
	Pkgs       []string `yaml:"pkgs,omitempty"`
}

type AptSpec struct {
	Name string `yaml:"name"`
}

type SnapSpec struct {
	Name string `yaml:"name"`
}

func (p *SnapSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return node.Decode(&p.Name)
	}
	type alias SnapSpec
	return node.Decode((*alias)(p))
}

type PipxSpec struct {
	Name string `yaml:"name"`
}

func (p *PipxSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return node.Decode(&p.Name)
	}
	type alias PipxSpec
	return node.Decode((*alias)(p))
}

type NPMSpec struct {
	Name string `yaml:"name"`
}

func (n *NPMSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return node.Decode(&n.Name)
	}
	type alias NPMSpec
	return node.Decode((*alias)(n))
}

type BrewTapSpec struct {
	Name string     `yaml:"name"`
	URL  *string    `yaml:"url,omitempty"`
	Pkgs []BrewSpec `yaml:"pkgs,omitempty"`
}

type BrewSpec struct {
	Name string `yaml:"name"`
	Cask bool   `yaml:"cask,omitempty"`
}

func (b *BrewSpec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		return node.Decode(&b.Name)
	}
	type alias BrewSpec
	return node.Decode((*alias)(b))
}

type GithubSpec struct {
	Name         string   `yaml:"name"`
	Repo         string   `yaml:"repo"`
	Version      string   `yaml:"version,omitempty"`
	Binaries     []string `yaml:"binaries,omitempty"`
	AssetPattern string   `yaml:"asset_pattern,omitempty"`
	BinaryPath   string   `yaml:"binary_path,omitempty"`
	InstallPath  string   `yaml:"install_path,omitempty"`
}

type BinarySpec struct {
	Name        string   `yaml:"name"`
	URL         string   `yaml:"url"`
	Version     string   `yaml:"version,omitempty"`
	BinaryPath  string   `yaml:"binary_path,omitempty"`
	InstallPath string   `yaml:"install_path,omitempty"`
	Binaries    []string `yaml:"binaries,omitempty"`
}

type ScriptSpec struct {
	Name   string `yaml:"name"`
	Script string `yaml:"script"`
}

type CustomSpec struct {
	Name         string `yaml:"name"`
	Install      string `yaml:"install"`
	Update       string `yaml:"update,omitempty"`
	InstallCheck string `yaml:"installCheck,omitempty"`
}
