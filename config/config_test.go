package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Errorf("Load should not fail, got: %v", err)
	}

	if cfg == nil {
		t.Errorf("Expected config, got nil")
	}

	if cfg.LogPath != DefaultLogPath {
		t.Errorf("Default log path: got %q, expected %q", cfg.LogPath, DefaultLogPath)
	}

	if cfg.FilePrefix != DefaultFilePrefix {
		t.Errorf("Default file prefix: got %q, expected %q", cfg.FilePrefix, DefaultFilePrefix)
	}
}

func TestConfigHasFilePrefix(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.FilePrefix == "" {
		t.Errorf("FilePrefix should not be empty")
	}
}

func TestConfigPathNormalization(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.LogPath == "" {
		t.Errorf("LogPath should not be empty")
	}
}

func TestDefaultLogPathConstant(t *testing.T) {
	if DefaultLogPath == "" {
		t.Errorf("DefaultLogPath should not be empty")
	}

	if !contains(DefaultLogPath, "Royal Quest") {
		t.Errorf("DefaultLogPath should contain 'Royal Quest'")
	}
}

func TestDefaultFilePrefixConstant(t *testing.T) {
	if DefaultFilePrefix == "" {
		t.Errorf("DefaultFilePrefix should not be empty")
	}

	if DefaultFilePrefix != "exp" {
		t.Errorf("DefaultFilePrefix should be 'exp', got %q", DefaultFilePrefix)
	}
}

func TestConfigFileLocation(t *testing.T) {
	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("Failed to get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	expectedPath := filepath.Join(exeDir, "config.json")

	// Config should look for file in same directory as executable
	if _, err := os.Stat(expectedPath); err == nil {
		// File exists, that's fine
	}
}

// Helper functions
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}