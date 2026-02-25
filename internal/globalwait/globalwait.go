// Package globalwait provides a global --wait mechanism for ionosctl.
// When --wait is set, it captures the href from the command's API response output,
// then polls that href until the resource reaches AVAILABLE state.
package globalwait

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"
)

var pollInterval = 5 * time.Second

const (
	httpTimeout = 10 * time.Second

	progressTpl = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".. " "..." "...." }}`
)

var (
	mu       sync.Mutex
	lastHref string
)

// CaptureHref extracts the href from the given API response data and stores it.
// It skips list responses (which have an "items" field) since those don't represent
// a single resource to wait on.
func CaptureHref(sourceData any) {
	href := ExtractHref(sourceData)
	if href == "" {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	lastHref = href
}

// GetHref returns the last captured href.
func GetHref() string {
	mu.Lock()
	defer mu.Unlock()
	return lastHref
}

// Reset clears the stored href. Call this before each command execution if needed.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	lastHref = ""
}

// ExtractHref extracts the "href" field from sourceData.
// Returns empty string if sourceData is a list (has "items" key), has no href,
// or cannot be parsed as JSON.
func ExtractHref(sourceData any) string {
	b, err := json.Marshal(sourceData)
	if err != nil {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return ""
	}
	// Skip list/collection responses
	if _, hasItems := m["items"]; hasItems {
		return ""
	}
	if href, ok := m["href"].(string); ok {
		return href
	}
	return ""
}

// WaitForAvailable polls the captured href until the resource reaches AVAILABLE state.
// Progress output is written to w (typically os.Stderr).
// Returns nil if no href was captured (command doesn't deal with API resources).
func WaitForAvailable(w io.Writer) error {
	href := GetHref()
	if href == "" {
		return nil
	}

	timeout := time.Duration(viper.GetInt(constants.ArgWaitTimeout)) * time.Second
	if timeout <= 0 {
		timeout = time.Duration(constants.DefaultWaitTimeoutSeconds) * time.Second
	}

	// Build full URL from href
	fullURL := buildFullURL(href)

	// Get auth credentials from the already-initialized client
	cl, err := client.Get()
	if err != nil {
		return fmt.Errorf("failed to get client for wait polling: %w", err)
	}
	cfg := cl.CloudClient.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Show progress bar on stderr
	bar := pb.New(1)
	bar.SetWriter(w)
	bar.SetTemplateString(progressTpl)
	bar.Start()
	defer bar.Finish()

	err = Poll(ctx, fullURL, cfg.Token, cfg.Username, cfg.Password)
	if err != nil {
		bar.SetTemplateString(progressTpl + " FAILED")
		return err
	}
	bar.SetTemplateString(progressTpl + " DONE")
	return nil
}

func buildFullURL(href string) string {
	// If href is already a full URL, use as-is
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return appendDepthParam(href)
	}

	// Otherwise, prepend the configured API base URL
	baseURL := viper.GetString(constants.ArgServerUrl)
	if baseURL == "" {
		baseURL = constants.DefaultApiURL
	}

	fullURL := strings.TrimRight(baseURL, "/") + href
	return appendDepthParam(fullURL)
}

func appendDepthParam(url string) string {
	if strings.Contains(url, "?") {
		return url + "&depth=1"
	}
	return url + "?depth=1"
}

// apiResponse is the minimal structure we need to parse from the polling response.
// We check both "state" (used by most services) and "status" (used by VPN).
type apiResponse struct {
	Metadata *apiMetadata `json:"metadata"`
}

type apiMetadata struct {
	State  string `json:"state"`
	Status string `json:"status"`
}

// Poll polls the given URL until the resource reaches a terminal ready state
// (AVAILABLE, ACTIVE, READY, DONE) or a failure state (FAILED).
func Poll(ctx context.Context, url, token, username, password string) error {
	httpClient := &http.Client{Timeout: httpTimeout}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	userAgent := viper.GetString(constants.CLIHttpUserAgent)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for resource to become AVAILABLE")
		case <-ticker.C:
		}

		state, err := fetchState(ctx, httpClient, url, token, username, password, userAgent)
		if err != nil {
			// Transient errors: retry on next tick
			continue
		}
		if state == "" {
			continue
		}

		switch strings.ToUpper(state) {
		case "AVAILABLE", "ACTIVE", "READY", "DONE":
			return nil
		case "FAILED":
			return fmt.Errorf("resource entered FAILED state")
		// BUSY, DEPLOYING, UPDATING, PROVISIONING, DESTROYING etc. - keep polling
		}
	}
}

func fetchState(ctx context.Context, httpClient *http.Client, url, token, username, password, userAgent string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else if username != "" {
		req.SetBasicAuth(username, password)
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	if body.Metadata == nil {
		return "", nil
	}

	// Most services use "state", VPN uses "status"
	state := body.Metadata.State
	if state == "" {
		state = body.Metadata.Status
	}
	return state, nil
}
