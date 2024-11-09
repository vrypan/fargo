package config

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	homeDir := func() string {
		t.Helper()
		usr, err := user.Current()
		if err != nil {
			t.Fatalf("failed to get current user: %v", err)
		}

		return usr.HomeDir
	}

	t.Run("XDG_CONFIG_HOME is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv(envXdgConfig, tmpDir)

		configDir, err := ConfigDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := filepath.Join(tmpDir, "fargo")
		if configDir != expected {
			if configDir == filepath.Join(homeDir(), ".fargo") {
				t.Fatal("you should run this test without ~/.fargo dir")
			}
			t.Fatalf("expected %v, got %v", expected, configDir)
		}
		if _, err := os.Stat(configDir); err != nil {
			t.Fatalf("failed to check config dir: %v", err)
		}
	})

	t.Run("XDG_CONFIG_HOME is not set", func(t *testing.T) {
		t.Setenv(envXdgConfig, "")

		configDir, err := ConfigDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := filepath.Join(homeDir(), ".fargo")
		if configDir != expected {
			t.Fatalf("expected %v, got %v", expected, configDir)
		}
		if _, err := os.Stat(configDir); err != nil {
			t.Fatalf("failed to check config dir: %v", err)
		}
	})
}
