package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"hermes/config"
	"hermes/hetzner"
	"hermes/ipcheck"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	dryRun := flag.Bool("dry-run", false, "Simulate actions without calling Hetzner API")
	configPath := flag.String("config", "/home/nixos/robot.json", "Path to the configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("hermes %s (commit: %s)\n", version, commit)
		return nil
	}

	// Initialize structured logger
	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	if *dryRun {
		slog.Info("DRY-RUN MODE ENABLED — No actual API calls will be made")
	}

	cfg, err := config.LoadFromPath(*configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	slog.Info("Starting Hetzner Failover Monitor", "failover_ip", cfg.FailoverIP)

	// Set up context with signal handling for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Check if the failover IP is present locally
	isLocal, err := ipcheck.IsIPLocal(cfg.FailoverIP)
	if err != nil {
		return fmt.Errorf("failed to check local IPs: %w", err)
	}

	if !isLocal {
		slog.Debug("Failover IP NOT found locally. No action taken.", "ip", cfg.FailoverIP)
		return nil
	}

	slog.Info("Failover IP detected locally. Ensuring routing...", "ip", cfg.FailoverIP)

	targetIP := cfg.MainIP
	if targetIP == "" {
		targetIP, err = ipcheck.GetMainIP()
		if err != nil {
			return fmt.Errorf("MainIP not set and failed to auto-detect: %w", err)
		}
		slog.Info("Auto-detected Main IP", "ip", targetIP)
	}

	client := hetzner.NewClient(cfg.HetznerUser, cfg.HetznerPass)
	if err := client.UpdateFailover(ctx, cfg.FailoverIP, targetIP, *dryRun); err != nil {
		return fmt.Errorf("failed to update Hetzner Failover: %w", err)
	}

	slog.Info("Successfully updated failover routing", "target_ip", targetIP)
	return nil
}
