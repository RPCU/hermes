// Package config handles the application configuration from files and environment variables.
package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Config holds all the configuration needed for Hermes.
type Config struct {
	HetznerUser string
	HetznerPass string
	FailoverIP  string
	MainIP      string
}

// Load loads the configuration from the default location and environment variables.
func Load() (*Config, error) {
	return LoadFromPath("/home/nixos/robot.json")
}

// LoadFromPath loads the configuration from a specific file path and environment variables.
func LoadFromPath(path string) (*Config, error) {
	cfg := &Config{}

	// Load from credentials file if it exists
	if f, err := os.ReadFile(path); err == nil {
		var fileCreds struct {
			User       string `json:"user"`
			Pass       string `json:"password"`
			FailoverIP string `json:"failover_ip"`
		}
		if err := json.Unmarshal(f, &fileCreds); err == nil {
			cfg.HetznerUser = fileCreds.User
			cfg.HetznerPass = fileCreds.Pass
			cfg.FailoverIP = fileCreds.FailoverIP
		}
	}

	// Environment variables override file
	if v := os.Getenv("HETZNER_USER"); v != "" {
		cfg.HetznerUser = v
	}
	if v := os.Getenv("HETZNER_PASS"); v != "" {
		cfg.HetznerPass = v
	}
	if v := os.Getenv("FAILOVER_IP"); v != "" {
		cfg.FailoverIP = v
	}
	if v := os.Getenv("MAIN_IP"); v != "" {
		cfg.MainIP = v
	}

	// Validate required fields
	if cfg.HetznerUser == "" {
		return nil, fmt.Errorf("HetznerUser is required (env HETZNER_USER or in %s)", path)
	}
	if cfg.HetznerPass == "" {
		return nil, fmt.Errorf("HetznerPass is required (env HETZNER_PASS or in %s)", path)
	}
	if cfg.FailoverIP == "" {
		return nil, fmt.Errorf("FailoverIP is required (env FAILOVER_IP or in %s)", path)
	}

	// Validate IP formats
	if net.ParseIP(cfg.FailoverIP) == nil {
		return nil, fmt.Errorf("invalid failover IP address: %s", cfg.FailoverIP)
	}
	if cfg.MainIP != "" && net.ParseIP(cfg.MainIP) == nil {
		return nil, fmt.Errorf("invalid main IP address: %s", cfg.MainIP)
	}

	return cfg, nil
}
