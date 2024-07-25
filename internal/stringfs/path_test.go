package stringfs

import (
	"os"
	"testing"
)

func TestParsePath(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("cant get current working dir:\n> " + err.Error())
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("cant get users home dir:\n> " + err.Error())
		return
	}

	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"", cwd, false},
		{"~", homeDir, false},
		{"~/test", homeDir + "/test", false},
		{"/absolute/path", "/absolute/path", false},
		{"/absolute/path/../test/test/../../path", "/absolute/path", false},
		{"relative/path", cwd + "/relative/path", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			err := ParsePath(&test.input)
			if test.expectError {
				if err == nil {
					t.Errorf("expected error for path %s, but got none", test.input)
					return
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for path %s, but got %v", test.input, err)
					return
				}
				if test.input != test.expected {
					t.Errorf("expected %s, but got %s", test.expected, test.input)
					return
				}
			}
		})
	}
}

func TestParsePathFrom(t *testing.T) {
	cwd := "/test/cwd"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("cant get users home dir:\n> " + err.Error())
		return
	}

	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"", cwd, false},
		{"~", homeDir, false},
		{"~/test", homeDir + "/test", false},
		{"/absolute/path", "/absolute/path", false},
		{"/absolute/path/../test/test/../../path", "/absolute/path", false},
		{"relative/path", cwd + "/relative/path", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			err := ParsePathFrom(&test.input, cwd)
			if test.expectError {
				if err == nil {
					t.Errorf("expected error for path %s, but got none", test.input)
					return
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for path %s, but got %v", test.input, err)
					return
				}
				if test.input != test.expected {
					t.Errorf("expected %s, but got %s", test.expected, test.input)
					return
				}
			}
		})
	}
}
