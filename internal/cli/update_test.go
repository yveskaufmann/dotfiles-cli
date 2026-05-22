package cli

import (
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name    string
		latest  string
		current string
		want    bool
	}{
		// dev builds always update
		{name: "dev build always updates", latest: "v1.0.0", current: "dev", want: true},

		// already up to date
		{name: "same version", latest: "v1.2.3", current: "v1.2.3", want: false},
		{name: "same version no v prefix on current", latest: "v1.2.3", current: "1.2.3", want: false},

		// patch bump
		{name: "newer patch", latest: "v1.2.4", current: "v1.2.3", want: true},
		{name: "older patch", latest: "v1.2.2", current: "v1.2.3", want: false},

		// minor bump
		{name: "newer minor", latest: "v1.3.0", current: "v1.2.9", want: true},
		{name: "older minor", latest: "v1.1.9", current: "v1.2.0", want: false},

		// major bump
		{name: "newer major", latest: "v2.0.0", current: "v1.9.9", want: true},
		{name: "older major", latest: "v1.0.0", current: "v2.0.0", want: false},

		// v prefix variants
		{name: "latest without v prefix", latest: "1.3.0", current: "v1.2.0", want: true},
		{name: "both without v prefix", latest: "1.3.0", current: "1.2.0", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNewerVersion(tt.latest, tt.current); got != tt.want {
				t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input string
		want  [3]int
	}{
		{"v1.2.3", [3]int{1, 2, 3}},
		{"1.2.3", [3]int{1, 2, 3}},
		{"v2.0.0", [3]int{2, 0, 0}},
		{"v0.0.1", [3]int{0, 0, 1}},
		{"v1.2.3-beta", [3]int{1, 2, 3}}, // pre-release suffix stripped
		{"v1.2", [3]int{1, 2, 0}},         // missing patch defaults to 0
		{"v1", [3]int{1, 0, 0}},            // missing minor/patch default to 0
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseVersion(tt.input); got != tt.want {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"v1.2.3", "v1.2.3"},
		{"1.2.3", "v1.2.3"},
		{"dev", "dev"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := normalizeVersion(tt.input); got != tt.want {
				t.Errorf("normalizeVersion(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
