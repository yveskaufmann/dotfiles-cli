package apt

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"yv35.com/dotfiles/internal/config"
	"yv35.com/dotfiles/internal/util/sh"
	"yv35.com/dotfiles/internal/util/stringutils"
)

type PPASourceTemplateData struct {
	URI        string
	Suites     string
	Components string
	Arch       string
	KeyFile    *string // pointer: nil means not set
}

func ReadPPASourceFromPPASpec(ppa config.PPASpec, keyFile string) (error, *PPASourceTemplateData) {

	// Default values
	suites := ppa.Suites
	if suites == "" {
		if out, err := sh.RunShellOutput("grep VERSION_CODENAME /etc/os-release | cut -d= -f2"); err == nil {
			suites = strings.TrimSpace(out)
		} else {
			suites = "stable"
		}
	}

	arch := "amd64"
	if out, err := sh.RunShellOutput("dpkg --print-architecture"); err == nil {
		arch = strings.TrimSpace(out)
	}

	finalURI := ppa.URI
	if finalURI == "" && strings.HasPrefix(ppa.Name, "ppa:") {
		ppaParts := strings.Split(strings.TrimPrefix(ppa.Name, "ppa:"), "/")
		if len(ppaParts) == 2 {
			finalURI = fmt.Sprintf("https://ppa.launchpadexternal.net/%s/%s/ubuntu", ppaParts[0], ppaParts[1])
		}
	}

	if finalURI == "" {
		return fmt.Errorf("PPA requires a URI or ppa:user/repo format got %s", ppa.URI), nil
	}

	var keyFilePtr *string
	if keyFile != "" {
		keyFilePtr = &keyFile
	}

	return nil, &PPASourceTemplateData{
		URI:        ppa.URI,
		Suites:     suites,
		Components: ppa.Components,
		Arch:       arch,
		KeyFile:    keyFilePtr,
	}
}

func RenderPPASourceFile(data *PPASourceTemplateData) string {

	if data.Components == "" {
		data.Components = "main"
	}

	if data.Suites == "" {
		data.Suites = "stable"
	}

	if data.Arch != "" {
		data.URI = strings.ReplaceAll(data.URI, ":arch", data.Arch)
	}

	tpl := stringutils.DepdupeIndention(` 
		Types: deb
		URIs: {{.URI}}
		Suites: {{.Suites}}
		Components: {{.Components}}
		Architectures: {{.Arch}}
		{{- if .KeyFile}}
		Signed-By: {{.KeyFile}}
		{{- end}}
	`)

	var buf bytes.Buffer
	template.Must(template.New("sourceFile").Parse(tpl)).Execute(&buf, data)
	return buf.String()
}

func WritePPASourceFile(sourcePath string, ppa config.PPASpec, sourceContent string) error {
	// Write to temp file and move
	tmpFile := "/tmp/" + ppa.SourceName + ".sources"
	if err := os.WriteFile(tmpFile, []byte(sourceContent), 0644); err != nil {
		return fmt.Errorf("failed to create temp sources file: %w", err)
	}

	slog.Info("Install source file", "from", tmpFile, "target", sourcePath)

	if err := sh.Run("sudo", "mv", tmpFile, sourcePath); err != nil {
		return fmt.Errorf("failed to move sources file: %w", err)
	}

	fmt.Printf("✅ Added PPA source: %s\n", ppa.Name)

	return nil
}
