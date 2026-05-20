package globalwait

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// poller holds shared HTTP client, auth, and user-agent for polling requests.
type poller struct {
	client    *http.Client
	token     string
	username  string
	password  string
	userAgent string
}

// provisioningFailure is returned by poll when a resource reaches a failure
// state (FAILED, ERROR) during provisioning. The API request itself succeeded,
// but the resource ended in a bad state. Callers print a warning instead of
// exiting with a non-zero code.
type provisioningFailure struct {
	url   string
	state string
}

func (e *provisioningFailure) Error() string {
	return fmt.Sprintf("resource at %s reached %s state", e.url, e.state)
}

func (p *poller) poll(ctx context.Context, url string, isDelete bool) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	firstSuccessfulPoll := true
	sawDestroying := false

	for {
		state, err := p.fetchState(ctx, url, isDelete || sawDestroying)
		if err != nil {
			firstSuccessfulPoll = false
			if strings.Contains(err.Error(), "authentication failed") ||
				strings.Contains(err.Error(), "server error") ||
				strings.Contains(err.Error(), "client error") {
				return err
			}
			// Other errors (network, bad JSON, 404 retrying) are transient, retry
		} else if state != "" {
			switch strings.ToUpper(state) {
			case "AVAILABLE", "ACTIVE", "READY", "DONE", "INACTIVE", "SUSPENDED":
				return nil
			case "FAILED", "ERROR":
				return &provisioningFailure{url: url, state: state}
			case "DESTROYING":
				// Resource being destroyed — keep polling until 404 returns "DONE".
				// This applies regardless of whether WE issued the delete;
				// e.g. "get --wait" on a resource deleted by another command.
				sawDestroying = true
			}
			firstSuccessfulPoll = false
		} else if firstSuccessfulPoll {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for %s to become AVAILABLE", url)
		case <-ticker.C:
		}
	}
}

type apiResponse struct {
	Metadata *apiMetadata `json:"metadata"`
}

type apiMetadata struct {
	State  string `json:"state"`
	Status string `json:"status"` // VPN uses "status" instead of "state"
}

func (p *poller) fetchState(ctx context.Context, url string, isDelete bool) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	p.setHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	// 404 during delete means resource is gone -> done.
	// 404 during create/update means resource is temporarily unavailable
	// during provisioning -> treat as transient, keep polling.
	if resp.StatusCode == http.StatusNotFound {
		if isDelete {
			return "DONE", nil
		}
		return "", fmt.Errorf("resource not found (HTTP 404), retrying")
	}

	// Auth errors should fail immediately, not retry for 10 minutes
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("authentication failed (HTTP %d) while polling resource state", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			if d, parseErr := strconv.Atoi(retryAfter); parseErr == nil && d > 0 {
				timer := time.NewTimer(time.Duration(d) * time.Second)
				select {
				case <-ctx.Done():
					timer.Stop()
					return "", ctx.Err()
				case <-timer.C:
				}
			}
		}
		return "", fmt.Errorf("rate limited (HTTP 429)")
	}

	// Other 4xx errors (400, 405, 409, etc.) are non-retryable client errors.
	// Without this, a 400 with valid JSON but no metadata would be treated as
	// "no state field" and the poller would silently declare success.
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return "", fmt.Errorf("client error (HTTP %d) while polling resource state", resp.StatusCode)
	}

	// Server errors are non-retryable. Retrying a 500 for 10 minutes wastes time.
	if resp.StatusCode >= 500 {
		return "", fmt.Errorf("server error (HTTP %d) while polling resource state", resp.StatusCode)
	}

	var body apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	if body.Metadata == nil {
		return "", nil
	}

	state := body.Metadata.State
	if state == "" {
		state = body.Metadata.Status
	}
	return state, nil
}

func (p *poller) fetchJSON(url string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	p.setHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to fetch resource (HTTP %d)", resp.StatusCode)
	}

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *poller) setHeaders(req *http.Request) {
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	} else if p.username != "" {
		req.SetBasicAuth(p.username, p.password)
	}
	if p.userAgent != "" {
		req.Header.Set("User-Agent", p.userAgent)
	}
	// Multi-contract users need this header; without it, polling requests
	// hit the default contract scope and may 404 on the wrong resource.
	if v, ok := os.LookupEnv("IONOS_CONTRACT_NUMBER"); ok && v != "" {
		req.Header.Set("X-Contract-Number", v)
	}
}
