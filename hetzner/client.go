package hetzner

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"io"
)

const baseURL = "https://robot-ws.your-server.de"

func UpdateFailover(user, pass, failoverIP, targetIP string, dryRun bool) error {
	endpoint := fmt.Sprintf("%s/failover/%s", baseURL, failoverIP)

	data := url.Values{}
	data.Set("active_server_ip", targetIP)

	if dryRun {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("🔍 DRY-RUN: Would execute the following:")
		fmt.Printf("  📍 Endpoint: %s\n", endpoint)
		fmt.Printf("  👤 User: %s\n", user)
		fmt.Printf("  🔑 Password: %s\n", maskPassword(pass))
		fmt.Printf("  📦 Data: active_server_ip=%s\n", targetIP)
		fmt.Println("")
		fmt.Println("📋 Equivalent curl command:")
		fmt.Printf("curl -X POST '%s' \\\n", endpoint)
		fmt.Printf("  -u '%s:%s' \\\n", user, maskPassword(pass))
		fmt.Printf("  -H 'Content-Type: application/x-www-form-urlencoded' \\\n")
		fmt.Printf("  -d 'active_server_ip=%s'\n", targetIP)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		return nil
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("Hetzner API Response: %s\n", string(bodyBytes))

	return nil
}

func maskPassword(pass string) string {
	if len(pass) <= 4 {
		return "****"
	}
	return pass[:2] + "****" + pass[len(pass)-2:]
}

