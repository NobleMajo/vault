package config

import (
	"os"
	"os/exec"
	"testing"
)

func TestParseConfigLockCommand(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"vault", "lock", "secret.txt"}
	cfg := ParseConfig("Demo", "demo", "1.0.0", "abc")

	if cfg.SubCommand != "lock" {
		t.Fatalf("SubCommand = %q, want lock", cfg.SubCommand)
	}
	if len(cfg.Args) != 1 || cfg.Args[0] != "secret.txt" {
		t.Fatalf("Args = %v, want [secret.txt]", cfg.Args)
	}
}

func TestParseConfigRejectsUnknownSubcommand(t *testing.T) {
	if os.Getenv("TEST_PARSE_CONFIG_SUBCMD") == "1" {
		os.Args = []string{"vault", "help"}
		ParseConfig("Demo", "demo", "1.0.0", "abc")
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestParseConfigRejectsUnknownSubcommand$")
	cmd.Env = append(os.Environ(), "TEST_PARSE_CONFIG_SUBCMD=1")
	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected exit error, got %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1", exitErr.ExitCode())
	}
}

func TestParseConfigPrintHelpExits(t *testing.T) {
	if os.Getenv("TEST_PARSE_CONFIG_PRINT_HELP") == "1" {
		os.Args = []string{"vault", "print", "--help"}
		ParseConfig("Demo", "demo", "1.0.0", "abc")
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestParseConfigPrintHelpExits$")
	cmd.Env = append(os.Environ(), "TEST_PARSE_CONFIG_PRINT_HELP=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
}

func TestParseConfigRootHelpExits(t *testing.T) {
	if os.Getenv("TEST_PARSE_CONFIG_ROOT_HELP") == "1" {
		os.Args = []string{"vault", "--help"}
		ParseConfig("Demo", "demo", "1.0.0", "abc")
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestParseConfigRootHelpExits$")
	cmd.Env = append(os.Environ(), "TEST_PARSE_CONFIG_ROOT_HELP=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
}

func TestParseConfigLockAlias(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"vault", "l", "secret.txt"}
	cfg := ParseConfig("Demo", "demo", "1.0.0", "abc")

	if cfg.SubCommand != "lock" {
		t.Fatalf("SubCommand = %q, want lock", cfg.SubCommand)
	}
}

func TestParseConfigUnlockCommand(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"vault", "unlock", "secret.vt"}
	cfg := ParseConfig("Demo", "demo", "1.0.0", "abc")

	if cfg.SubCommand != "unlock" {
		t.Fatalf("SubCommand = %q, want unlock", cfg.SubCommand)
	}
}

func TestParseConfigTempSecondsFlag(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"vault", "temp", "secret.vt", "-t", "30"}
	cfg := ParseConfig("Demo", "demo", "1.0.0", "abc")

	if cfg.SubCommand != "temp" {
		t.Fatalf("SubCommand = %q, want temp", cfg.SubCommand)
	}
	if cfg.TempDecodeSeconds != 30 {
		t.Fatalf("TempDecodeSeconds = %d, want 30", cfg.TempDecodeSeconds)
	}
}

func TestParseConfigVerboseEnv(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	t.Setenv("VAULT_VERBOSE", "true")
	os.Args = []string{"vault", "lock", "secret.txt"}
	cfg := ParseConfig("Demo", "demo", "1.0.0", "abc")

	if !cfg.Verbose {
		t.Fatal("expected VAULT_VERBOSE to enable verbose mode")
	}
}

func TestParseConfigVaultExtEnv(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	t.Setenv("VAULT_EXT", "enc")
	os.Args = []string{"vault", "lock", "secret.txt"}
	cfg := ParseConfig("Demo", "demo", "1.0.0", "abc")

	if cfg.VaultFileExtension != "enc" {
		t.Fatalf("VaultFileExtension = %q, want enc", cfg.VaultFileExtension)
	}
}

func TestParseConfigVersionSubcommand(t *testing.T) {
	if os.Getenv("TEST_PARSE_CONFIG_VERSION") == "1" {
		os.Args = []string{"vault", "version"}
		ParseConfig("Demo", "demo", "1.0.0", "abc")
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestParseConfigVersionSubcommand$")
	cmd.Env = append(os.Environ(), "TEST_PARSE_CONFIG_VERSION=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
}

func TestParseConfigBareVaultShowsHelpExits(t *testing.T) {
	if os.Getenv("TEST_PARSE_CONFIG_BARE") == "1" {
		os.Args = []string{"vault"}
		ParseConfig("Demo", "demo", "1.0.0", "abc")
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestParseConfigBareVaultShowsHelpExits$")
	cmd.Env = append(os.Environ(), "TEST_PARSE_CONFIG_BARE=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got %v", err)
	}
}
