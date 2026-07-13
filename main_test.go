package main

import (
	"testing"

	"github.com/NobleMajo/vault/internal/config"
)

func TestTargetFileDefault(t *testing.T) {
	got := targetFile(&config.AppConfig{})
	if got != "vault" {
		t.Fatalf("targetFile() = %q, want vault", got)
	}
}

func TestTargetFileStripsVaultExtension(t *testing.T) {
	got := targetFile(&config.AppConfig{
		Args:               []string{"secret.vt"},
		VaultFileExtension: "vt",
		PlainFileExtension: "txt",
	})
	if got != "secret" {
		t.Fatalf("targetFile() = %q, want secret", got)
	}
}

func TestTargetFileStripsPlainExtension(t *testing.T) {
	got := targetFile(&config.AppConfig{
		Args:               []string{"notes.txt"},
		VaultFileExtension: "vt",
		PlainFileExtension: "txt",
	})
	if got != "notes" {
		t.Fatalf("targetFile() = %q, want notes", got)
	}
}

func TestTargetFileKeepsBaseNameWithoutKnownExtension(t *testing.T) {
	got := targetFile(&config.AppConfig{
		Args:               []string{"secret"},
		VaultFileExtension: "vt",
		PlainFileExtension: "txt",
	})
	if got != "secret" {
		t.Fatalf("targetFile() = %q, want secret", got)
	}
}
