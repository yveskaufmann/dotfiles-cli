package nvm

import (
	"testing"

	"yv35.com/dotfiles-cli/internal/config"
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
			name: "has nvm versions",
			group: config.DependencyGroup{
				Name: "test-group",
				NVM: []config.NVMSpec{
					{Default: "18", Versions: []string{"20"}},
				},
			},
			want: true,
		},
		{
			name: "has single nvm version",
			group: config.DependencyGroup{
				Name: "test-group",
				NVM: []config.NVMSpec{
					{Default: "18"},
				},
			},
			want: true,
		},
		{
			name: "has other provider fields only",
			group: config.DependencyGroup{
				Name: "test-group",
				Apt:  []string{"vim"},
				NPM: []config.NPMSpec{
					{Name: "typescript"},
				},
			},
			want: false,
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
