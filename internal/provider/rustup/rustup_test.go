package rustup

import (
	"fmt"
	"testing"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
	"yv35.com/dotfiles-cli/internal/util/sh"
)

func TestProvider_HasConfig(t *testing.T) {
	tests := []struct {
		name  string
		group config.DependencyGroup
		want  bool
	}{
		{
			name: "empty group",
			group: config.DependencyGroup{
				Name: "test-group",
			},
			want: false,
		},
		{
			name: "has rustup config",
			group: config.DependencyGroup{
				Name: "test-group",
				RustUp: []config.RustUpSpec{
					{DefaultToolchain: "stable"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProvider()
			if got := p.HasConfig(tt.group); got != tt.want {
				t.Errorf("Provider.HasConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_Setup(t *testing.T) {
	p := NewProvider()

	if err := p.Setup(); err != nil {
		t.Errorf("Provider.Setup() error = %v", err)
	}

	// Check if rustup is installed
	if err := sh.RunShell("rustup --version > /dev/null 2>&1"); err != nil {
		t.Errorf("Rustup not found on PATH: %v", err)
	}
}

func TestProvider_InstallComponents(t *testing.T) {
	p := NewProvider()
	p.Setup()

	onComplete := func(result types.TaskResult) {
		if result.Error != nil {
			t.Errorf("Task %s failed: %v", result.Name, result.Error)
		}

		if result.Status != types.StatusSuccess {
			t.Errorf("Task %s did not succeed: status %v", result.Name, result.Status)
		}

		fmt.Printf("Task %s completed\n", result.Name)
	}

	p.Install(config.DependencyGroup{
		RustUp: []config.RustUpSpec{
			{
				Components: []string{"clobby", "rustfmt"},
			},
		},
	}, onComplete)

}
