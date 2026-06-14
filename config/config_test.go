package config

import (
	"os"
	"testing"
)

func TestLoadFromPath(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "robot.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{"user": "testuser", "password": "testpass", "failover_ip": "1.2.3.4"}`
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	t.Run("Load from file", func(t *testing.T) {
		cfg, err := LoadFromPath(tmpFile.Name())
		if err != nil {
			t.Fatalf("LoadFromPath failed: %v", err)
		}

		if cfg.HetznerUser != "testuser" {
			t.Errorf("expected testuser, got %s", cfg.HetznerUser)
		}
		if cfg.HetznerPass != "testpass" {
			t.Errorf("expected testpass, got %s", cfg.HetznerPass)
		}
		if cfg.FailoverIP != "1.2.3.4" {
			t.Errorf("expected 1.2.3.4, got %s", cfg.FailoverIP)
		}
	})

	t.Run("Override with environment variables", func(t *testing.T) {
		os.Setenv("HETZNER_USER", "envuser")
		defer os.Unsetenv("HETZNER_USER")

		cfg, err := LoadFromPath(tmpFile.Name())
		if err != nil {
			t.Fatalf("LoadFromPath failed: %v", err)
		}

		if cfg.HetznerUser != "envuser" {
			t.Errorf("expected envuser, got %s", cfg.HetznerUser)
		}
		// Password should still come from file
		if cfg.HetznerPass != "testpass" {
			t.Errorf("expected testpass, got %s", cfg.HetznerPass)
		}
	})

	t.Run("Missing required fields", func(t *testing.T) {
		emptyFile, err := os.CreateTemp("", "empty.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(emptyFile.Name())
		if _, err := emptyFile.Write([]byte(`{}`)); err != nil {
			t.Fatal(err)
		}
		emptyFile.Close()

		_, err = LoadFromPath(emptyFile.Name())
		if err == nil {
			t.Error("expected error for missing required fields, got nil")
		}
	})

	t.Run("Invalid JSON syntax", func(t *testing.T) {
		badFile, err := os.CreateTemp("", "bad.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(badFile.Name())
		if _, err := badFile.Write([]byte(`{not valid json`)); err != nil {
			t.Fatal(err)
		}
		badFile.Close()

		// Should still fail because required fields are missing
		// (invalid JSON is silently ignored, falling back to env vars)
		_, err = LoadFromPath(badFile.Name())
		if err == nil {
			t.Error("expected error for invalid JSON with no env vars, got nil")
		}
	})

	t.Run("Invalid failover IP format", func(t *testing.T) {
		badIPFile, err := os.CreateTemp("", "badip.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(badIPFile.Name())
		if _, err := badIPFile.Write([]byte(`{"user": "u", "password": "p", "failover_ip": "not-an-ip"}`)); err != nil {
			t.Fatal(err)
		}
		badIPFile.Close()

		_, err = LoadFromPath(badIPFile.Name())
		if err == nil {
			t.Error("expected error for invalid IP, got nil")
		}
	})
}
