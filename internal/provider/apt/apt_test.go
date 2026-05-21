package apt

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
			name: "has apt packages",
			group: config.DependencyGroup{
				Name: "test-group",
				Apt:  []string{"vim", "git"},
			},
			want: true,
		},
		{
			name: "has single apt package",
			group: config.DependencyGroup{
				Name: "test-group",
				Apt:  []string{"vim"},
			},
			want: true,
		},
		{
			name: "has PPA only",
			group: config.DependencyGroup{
				Name: "test-group",
				PPA: []config.PPASpec{
					{Name: "ppa:test/repo"},
				},
			},
			want: true,
		},
		{
			name: "has both apt and PPA",
			group: config.DependencyGroup{
				Name: "test-group",
				Apt:  []string{"vim"},
				PPA: []config.PPASpec{
					{Name: "ppa:test/repo"},
				},
			},
			want: true,
		},
		{
			name: "has other provider fields only",
			group: config.DependencyGroup{
				Name: "test-group",
				Brew: []config.BrewSpec{
					{Name: "vim"},
				},
				Snap: []config.SnapSpec{
					{Name: "code"},
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
