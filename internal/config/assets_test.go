package config

import "testing"

func TestEnvIsString(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		t.Setenv("TEST_STRING", "hello")
		var got string
		EnvIsString("TEST_STRING", func(value string) {
			got = value
		})
		if got != "hello" {
			t.Fatalf("got %q, want hello", got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv("TEST_STRING", "")
		called := false
		EnvIsString("TEST_STRING", func(string) { called = true })
		if called {
			t.Fatal("expected callback to be skipped for empty env")
		}
	})
}

func TestEnvIsInt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Setenv("TEST_INT", "42")
		var got int
		EnvIsInt("TEST_INT", func(value int) { got = value })
		if got != 42 {
			t.Fatalf("got %d, want 42", got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv("TEST_INT", "nope")
		called := false
		EnvIsInt("TEST_INT", func(int) { called = true })
		if called {
			t.Fatal("expected callback to be skipped for invalid int")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv("TEST_INT", "")
		called := false
		EnvIsInt("TEST_INT", func(int) { called = true })
		if called {
			t.Fatal("expected callback to be skipped for empty env")
		}
	})
}

func TestEnvIsBool(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Setenv("TEST_BOOL", "true")
		var got bool
		EnvIsBool("TEST_BOOL", func(value bool) { got = value })
		if !got {
			t.Fatal("expected true")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Setenv("TEST_BOOL", "maybe")
		called := false
		EnvIsBool("TEST_BOOL", func(bool) { called = true })
		if called {
			t.Fatal("expected callback to be skipped for invalid bool")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv("TEST_BOOL", "")
		called := false
		EnvIsBool("TEST_BOOL", func(bool) { called = true })
		if called {
			t.Fatal("expected callback to be skipped for empty env")
		}
	})
}
