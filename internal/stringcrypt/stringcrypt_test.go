package stringcrypt

import (
	"os"
	"path/filepath"
	"testing"
)

// Assuming LoadRsaPublicKey and LoadRsaPrivateKey are defined in the same package

func TestLoadRsaPublicKey(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("cant get current working dir:\n> " + err.Error())
	}

	tests := []struct {
		path      string
		expectErr bool
	}{
		{wd + "/../../keys/test_id_rsa.pub.pem", false},
		{wd + "/../../keys/test_id_rsa.pub", false},
		{wd + "/../../keys/nonexistent.pub", true},
	}

	for _, test := range tests {
		_, file := filepath.Split(test.path)
		t.Run(file, func(t *testing.T) {
			key, err := LoadRsaPublicKey(test.path)
			if test.expectErr {
				if err == nil {
					t.Errorf("expected error for key %s, but got none", file)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for key %s, but got:\n"+err.Error(), file)
				} else if key == nil {
					t.Errorf("expected rsa.PublicKey, but got nil")
				}
			}
		})
	}
}

func TestLoadRsaPrivateKey(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("cant get current working dir:\n> " + err.Error())
	}

	tests := []struct {
		path      string
		expectErr bool
	}{
		{wd + "/../../keys/test_id_rsa", false},
		{wd + "/../../keys/test_id_rsa.rsa", false},
		{wd + "/../../keys/nonexistent.rsa", true},
	}

	for _, test := range tests {
		_, file := filepath.Split(test.path)
		t.Run(file, func(t *testing.T) {
			key, err := LoadRsaPrivateKey(test.path)
			if test.expectErr {
				if err == nil {
					t.Errorf("expected error for key %s, but got none", file)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for key %s, but got:\n"+err.Error(), file)
				} else if key == nil {
					t.Errorf("expected rsa.PrivateKey, but got nil")
				}
			}
		})
	}
}
