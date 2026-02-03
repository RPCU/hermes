package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HetznerUser string
	HetznerPass string
	FailoverIP  string
	MainIP      string
}

func Load() (*Config, error) {
	return LoadWithDryRun(false)
}

func LoadWithDryRun(dryRun bool) (*Config, error) {
	cfg := &Config{}

	// Load from credentials file
	credsPath := "/home/nixos/robot.json"
	if f, err := os.ReadFile(credsPath); err == nil {
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
		return nil, fmt.Errorf("HetznerUser is required (env HETZNER_USER or %s)", credsPath)
	}
	if cfg.HetznerPass == "" {
		return nil, fmt.Errorf("HetznerPass is required (env HETZNER_PASS or %s)", credsPath)
	}
	if cfg.FailoverIP == "" {
		return nil, fmt.Errorf("FailoverIP is required (env FAILOVER_IP or %s)", credsPath)
	}

	return cfg, nil
}
