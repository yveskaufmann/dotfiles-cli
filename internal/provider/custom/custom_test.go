package custom

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
			name: "has custom packages",
			group: config.DependencyGroup{
				Name: "test-group",
				Custom: []config.CustomSpec{
					{Name: "custom1", Install: "echo 'install custom1'"},
					{Name: "custom2", Install: "echo 'install custom2'"},
				},
			},
			want: true,
		},
		{
			name: "has single custom package",
			group: config.DependencyGroup{
				Name: "test-group",
				Custom: []config.CustomSpec{
					{Name: "custom1", Install: "echo 'install custom1'"},
				},
			},
			want: true,
		},
		{
			name: "has other provider fields only",
			group: config.DependencyGroup{
				Name: "test-group",
				Apt:  []string{"vim"},
				Brew: []config.BrewSpec{
					{Name: "vim"},
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
