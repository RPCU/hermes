package ipcheck

import (
	"fmt"
	"net"
)

// IsIPLocal checks if the given IP address exists on any local interface.
func IsIPLocal(ipStr string) (bool, error) {
	targetIP := net.ParseIP(ipStr)
	if targetIP == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return false, fmt.Errorf("failed to list interfaces: %w", err)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.Equal(targetIP) {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetMainIP determines the outbound IP of this machine by dialing a public address.
func GetMainIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
