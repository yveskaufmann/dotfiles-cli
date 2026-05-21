package archive

import (
	"os"
	"path/filepath"
	"testing"
)

// --- isWritable ---

func TestIsWritable_WritableExistingFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "binary")
	if err := os.WriteFile(dest, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	if !isWritable(dir, dest) {
		t.Error("expected writable file to return true")
	}
}

func TestIsWritable_NonWritableExistingFile(t *testing.T) {
	// Reproduces the crane scenario: binary exists inside $HOME but was
	// installed previously with root permissions (mode 0444 simulates that).
	dir := t.TempDir()
	dest := filepath.Join(dir, "crane")
	if err := os.WriteFile(dest, []byte("old binary"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dest, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dest, 0644) // restore so t.TempDir cleanup can delete the file

	if isWritable(dir, dest) {
		t.Error("expected non-writable existing file to return false")
	}
}

func TestIsWritable_WritableDir_DestAbsent(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "new-binary")

	if !isWritable(dir, dest) {
		t.Error("expected writable directory with absent dest to return true")
	}
}

func TestIsWritable_NonWritableDir_DestAbsent(t *testing.T) {
	// Reproduces the root-owned ~/bin scenario.
	parent := t.TempDir()
	dir := filepath.Join(parent, "readonly-bin")
	if err := os.MkdirAll(dir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0755) // restore so t.TempDir cleanup can remove it

	dest := filepath.Join(dir, "new-binary")

	if isWritable(dir, dest) {
		t.Error("expected non-writable directory to return false")
	}
}

func TestIsWritable_DirAbsent_WritableParent(t *testing.T) {
	// Neither targetDir nor dest exist yet; parent is writable.
	parent := t.TempDir()
	dir := filepath.Join(parent, "bin")
	dest := filepath.Join(dir, "binary")

	if !isWritable(dir, dest) {
		t.Error("expected missing dir with writable parent to return true")
	}
}

// --- InstallBinary ---

func TestInstallBinary_WritableTarget(t *testing.T) {
	// targetDir must be inside $HOME so the path-prefix heuristic does not
	// trigger sudo (which requires a terminal and would fail in CI).
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	targetDir, err := os.MkdirTemp(home, ".dotfiles_test_install_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(targetDir)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "mytool")
	if err := os.WriteFile(src, []byte("#!/bin/sh\necho hello"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := InstallBinary(src, targetDir, "mytool"); err != nil {
		t.Fatalf("InstallBinary() error = %v", err)
	}

	dest := filepath.Join(targetDir, "mytool")
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatalf("binary not found at %s: %v", dest, err)
	}
	if info.Mode()&0111 == 0 {
		t.Error("installed binary is not executable")
	}
}

func TestInstallBinary_NonWritableExistingFile_DetectedBeforeMove(t *testing.T) {
	// Verifies that isWritable correctly detects a non-writable pre-existing
	// binary (the crane bug). We assert on isWritable directly because
	// invoking sudo in a test environment is not feasible.
	dir := t.TempDir()
	dest := filepath.Join(dir, "crane")
	if err := os.WriteFile(dest, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dest, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dest, 0644)

	if isWritable(dir, dest) {
		t.Error("isWritable should return false for a non-writable existing file, " +
			"causing InstallBinary to escalate to sudo")
	}
}
