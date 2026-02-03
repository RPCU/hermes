package main

import (
	"flag"
	"log"

	"hermes/config"
	"hermes/hetzner"
	"hermes/ipcheck"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Simulate actions without calling Hetzner API")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if *dryRun {
		log.Println("🔍 DRY-RUN MODE ENABLED - No actual API calls will be made")
	}

	cfg, err := config.LoadWithDryRun(*dryRun)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Hetzner Failover Monitor for IP: %s", cfg.FailoverIP)

	// Check if the failover IP is present locally
	isLocal, err := ipcheck.IsIPLocal(cfg.FailoverIP)
	if err != nil {
		log.Fatalf("Failed to check local IPs: %v", err)
	}

	if isLocal {
		log.Printf("Failover IP %s detected locally. Ensuring routing...", cfg.FailoverIP)

		targetIP := cfg.MainIP
		if targetIP == "" {
			var err error
			targetIP, err = ipcheck.GetMainIP()
			if err != nil {
				log.Fatalf("MainIP not set and failed to auto-detect: %v", err)
			}
			log.Printf("Auto-detected Main IP: %s", targetIP)
		}

		err := hetzner.UpdateFailover(cfg.HetznerUser, cfg.HetznerPass, cfg.FailoverIP, targetIP, *dryRun)
		if err != nil {
			log.Fatalf("Failed to update Hetzner Failover: %v", err)
		}

		log.Printf("Successfully updated failover routing to %s", targetIP)
	} else {
		log.Printf("Failover IP %s NOT found locally. No action taken.", cfg.FailoverIP)
	}
}
