package stringfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.txt")

	if err := SafeWriteFile(path, "payload", 0o600); err != nil {
		t.Fatalf("SafeWriteFile: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got != "payload" {
		t.Fatalf("ReadFile = %q, want payload", got)
	}

	tmpPath := filepath.Join(dir, ".tmp_secret.txt")
	if Exists(tmpPath) {
		t.Fatalf("expected temp file removed after rename, still exists at %s", tmpPath)
	}
}

func TestRemoveTmpSafeFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.txt")
	tmpPath := filepath.Join(dir, ".tmp_secret.txt")

	if err := os.WriteFile(tmpPath, []byte("tmp"), 0o600); err != nil {
		t.Fatalf("WriteFile tmp: %v", err)
	}

	RemoveTmpSafeFile(path)
	if Exists(tmpPath) {
		t.Fatal("expected temp file removed")
	}
}
