package jetbrains

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
			name:  "empty group",
			group: config.DependencyGroup{Name: "test-group"},
			want:  false,
		},
		{
			name: "has jetbrains specs",
			group: config.DependencyGroup{
				Name: "test-group",
				Jetbrains: []config.JetbrainsSpec{
					{IDE: "IIU"},
					{IDE: "WS"},
				},
			},
			want: true,
		},
		{
			name: "has single jetbrains spec",
			group: config.DependencyGroup{
				Name:      "test-group",
				Jetbrains: []config.JetbrainsSpec{{IDE: "PCP", Version: "2024.1.1"}},
			},
			want: true,
		},
		{
			name: "has other provider fields only",
			group: config.DependencyGroup{
				Name: "test-group",
				Apt:  []string{"vim"},
				Brew: []config.BrewSpec{{Name: "vim"}},
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

func TestResolveCode(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"IIU", "IIU", false},
		{"WS", "WS", false},
		{"idea-IU", "IIU", false},
		{"IntelliJ IDEA Ultimate", "IIU", false},
		{"WebStorm", "WS", false},
		{"unknown-ide", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := resolveCode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveCode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resolveCode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestInstallDir(t *testing.T) {
	tests := []struct {
		code    string
		wantErr bool
	}{
		{"IIU", false},
		{"WS", false},
		{"PCP", false},
		{"CL", false},
		{"UNKNOWN", true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			_, err := installDir(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("installDir(%q) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			}
		})
	}
}
