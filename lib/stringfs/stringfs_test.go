package stringfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReadRemoveFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := WriteFile(path, "hello", 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got != "hello" {
		t.Fatalf("ReadFile = %q, want hello", got)
	}

	if !Exists(path) {
		t.Fatal("expected file to exist")
	}
	exists, isFile := IsFile(path)
	if !exists || !isFile {
		t.Fatalf("IsFile = (%v, %v), want (true, true)", exists, isFile)
	}
	exists, isDir := IsDir(path)
	if !exists || isDir {
		t.Fatalf("IsDir = (%v, %v), want (true, false)", exists, isDir)
	}

	if err := RemoveFile(path); err != nil {
		t.Fatalf("RemoveFile: %v", err)
	}
	if Exists(path) {
		t.Fatal("expected file removed")
	}
}

func TestRemoveFileMissingPathIsNoOp(t *testing.T) {
	if err := RemoveFile(filepath.Join(t.TempDir(), "missing.txt")); err != nil {
		t.Fatalf("RemoveFile missing: %v", err)
	}
}

func TestRemoveFileRemovesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	if err := RemoveFile(dir); err != nil {
		t.Fatalf("RemoveFile dir: %v", err)
	}
	if Exists(dir) {
		t.Fatal("expected directory removed")
	}
}
