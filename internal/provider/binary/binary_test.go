package binary

import (
	"testing"

	"yv35.com/dotfiles/internal/config"
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
			name: "has binary specs",
			group: config.DependencyGroup{
				Name: "test-group",
				Binary: []config.BinarySpec{
					{Name: "tool1", URL: "https://example.com/tool1"},
					{Name: "tool2", URL: "https://example.com/tool2"},
				},
			},
			want: true,
		},
		{
			name: "has single binary spec",
			group: config.DependencyGroup{
				Name: "test-group",
				Binary: []config.BinarySpec{
					{Name: "tool1", URL: "https://example.com/tool1"},
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
