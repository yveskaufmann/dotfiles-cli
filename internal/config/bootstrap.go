package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Repository represents a git repository configuration
type Repository struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"` // ssh or https
}

// BootstrapConfig represents the bootstrap configuration
type BootstrapConfig struct {
	Dotfiles struct {
		Repository Repository `yaml:"repository"`
	} `yaml:"dotfiles"`
}

// LoadBootstrapConfig loads the bootstrap configuration from ~/.config/.dotfiles/config.yaml
// Returns nil if the file doesn't exist (first-time bootstrap scenario)
func LoadBootstrapConfig() (*BootstrapConfig, error) {
	configPath := getConfigPath()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil // No config file, first-time bootstrap
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config BootstrapConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveBootstrapConfig saves the bootstrap configuration to ~/.config/.dotfiles/config.yaml
func SaveBootstrapConfig(config *BootstrapConfig) error {
	configPath := getConfigPath()

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath returns the path to the bootstrap config file
func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", ".dotfiles", "config.yaml")
}
