// Package hetzner provides a client for interacting with the Hetzner Robot API.
package hetzner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://robot-ws.your-server.de"

// Client is a Hetzner Robot API client.
type Client struct {
	BaseURL    string
	User       string
	Pass       string
	HTTPClient *http.Client
}

// NewClient creates a new Hetzner Robot API client with sensible defaults.
func NewClient(user, pass string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		User:    user,
		Pass:    pass,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UpdateFailover updates the routing of a failover IP to a target IP.
func (c *Client) UpdateFailover(ctx context.Context, failoverIP, targetIP string, dryRun bool) error {
	endpoint := fmt.Sprintf("%s/failover/%s", c.BaseURL, failoverIP)

	data := url.Values{}
	data.Set("active_server_ip", targetIP)

	if dryRun {
		slog.Info(
			"DRY-RUN: simulating API call",
			"endpoint", endpoint,
			"user", c.User,
			"active_server_ip", targetIP,
		)
		slog.Debug(
			"DRY-RUN: equivalent curl command",
			"command", fmt.Sprintf(
				"curl -X POST '%s' -u '%s:PASSWORD' -H 'Content-Type: application/x-www-form-urlencoded' -d 'active_server_ip=%s'",
				endpoint, c.User, targetIP,
			),
		)
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.User, c.Pass)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Debug("API response", "status", resp.StatusCode, "body", string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
