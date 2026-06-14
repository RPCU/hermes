package hetzner

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateFailover(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/failover/1.2.3.4" {
				t.Errorf("expected /failover/1.2.3.4, got %s", r.URL.Path)
			}

			user, pass, ok := r.BasicAuth()
			if !ok || user != "testuser" || pass != "testpass" {
				t.Errorf("basic auth failed: ok=%v user=%s", ok, user)
			}

			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("unexpected content type: %s", r.Header.Get("Content-Type"))
			}

			body, _ := io.ReadAll(r.Body)
			if string(body) != "active_server_ip=5.6.7.8" {
				t.Errorf("unexpected body: %s", string(body))
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"failover": {"ip": "1.2.3.4", "active_server_ip": "5.6.7.8"}}`)
		}))
		defer server.Close()

		client := NewClient("testuser", "testpass")
		client.BaseURL = server.URL

		err := client.UpdateFailover(ctx, "1.2.3.4", "5.6.7.8", false)
		if err != nil {
			t.Fatalf("UpdateFailover failed: %v", err)
		}
	})

	t.Run("API error returns descriptive error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": {"status": 400, "message": "invalid request"}}`)
		}))
		defer server.Close()

		client := NewClient("testuser", "testpass")
		client.BaseURL = server.URL

		err := client.UpdateFailover(ctx, "1.2.3.4", "5.6.7.8", false)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Dry run skips API call", func(t *testing.T) {
		client := NewClient("testuser", "testpass")
		client.BaseURL = "http://invalid-should-not-be-called"

		err := client.UpdateFailover(ctx, "1.2.3.4", "5.6.7.8", true)
		if err != nil {
			t.Fatalf("Dry run failed: %v", err)
		}
	})

	t.Run("Cancelled context returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // cancel immediately

		client := NewClient("testuser", "testpass")
		client.BaseURL = server.URL

		err := client.UpdateFailover(cancelledCtx, "1.2.3.4", "5.6.7.8", false)
		if err == nil {
			t.Error("expected error for cancelled context, got nil")
		}
	})

	t.Run("Status 201 is accepted", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"failover": {"ip": "1.2.3.4"}}`)
		}))
		defer server.Close()

		client := NewClient("testuser", "testpass")
		client.BaseURL = server.URL

		err := client.UpdateFailover(ctx, "1.2.3.4", "5.6.7.8", false)
		if err != nil {
			t.Fatalf("expected success for 201, got: %v", err)
		}
	})
}

func TestNewClient(t *testing.T) {
	client := NewClient("user", "pass")

	if client.BaseURL != DefaultBaseURL {
		t.Errorf("expected base URL %s, got %s", DefaultBaseURL, client.BaseURL)
	}
	if client.User != "user" {
		t.Errorf("expected user 'user', got '%s'", client.User)
	}
	if client.HTTPClient == nil {
		t.Error("expected non-nil HTTP client")
	}
	if client.HTTPClient.Timeout == 0 {
		t.Error("expected non-zero HTTP timeout")
	}
}
