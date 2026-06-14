package ipcheck

import (
	"testing"
)

func TestIsIPLocal(t *testing.T) {
	// 127.0.0.1 should be local on almost any system
	isLocal, err := IsIPLocal("127.0.0.1")
	if err != nil {
		t.Fatalf("IsIPLocal failed: %v", err)
	}
	if !isLocal {
		t.Errorf("expected 127.0.0.1 to be local")
	}

	// 255.255.255.255 should NOT be local
	isLocal, err = IsIPLocal("255.255.255.255")
	if err != nil {
		t.Fatalf("IsIPLocal failed: %v", err)
	}
	if isLocal {
		t.Errorf("expected 255.255.255.255 to not be local")
	}

	// Invalid IP
	_, err = IsIPLocal("invalid-ip")
	if err == nil {
		t.Error("expected error for invalid IP, got nil")
	}
}

func TestGetMainIP(t *testing.T) {
	// This test requires internet access, but we can at least check if it returns a valid IP
	ip, err := GetMainIP()
	if err != nil {
		t.Skipf("Skipping GetMainIP test (likely no internet): %v", err)
	}
	if ip == "" {
		t.Error("expected non-empty IP")
	}
}
