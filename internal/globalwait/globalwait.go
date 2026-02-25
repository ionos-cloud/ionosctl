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

// RenderInfo stores the parameters needed to re-render output after waiting.
// This is captured from GenerateOutput so we can re-render with fresh data.
// Deprecated: Use CaptureRerenderable with the table package instead.
type RenderInfo struct {
	Prefix  string
	Mapping map[string]string
	Cols    []string
}

// Rerenderable can re-render its output with fresh source data.
// Implemented by the table.Table type for seamless --wait integration.
type Rerenderable interface {
	Extract(sourceData any) error
	Render(visibleCols []string) (string, error)
}

var (
	mu               sync.Mutex
	lastHref         string
	lastRenderInfo   *RenderInfo   // Legacy: for jsontabwriter
	lastRerenderable Rerenderable  // New: for table package
	lastVisibleCols  []string      // New: cols for rerenderable
	rerendering      bool
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

// CaptureRenderInfo stores the rendering parameters so output can be re-rendered
// with fresh data after waiting completes.
func CaptureRenderInfo(prefix string, mapping map[string]string, cols []string) {
	mu.Lock()
	defer mu.Unlock()
	lastRenderInfo = &RenderInfo{Prefix: prefix, Mapping: mapping, Cols: cols}
}

// GetRenderInfo returns the captured render info, or nil if not captured.
func GetRenderInfo() *RenderInfo {
	mu.Lock()
	defer mu.Unlock()
	return lastRenderInfo
}

// IsRerendering returns true when we're in the re-render pass after waiting.
// This prevents GenerateOutput from suppressing output during the second call.
func IsRerendering() bool {
	mu.Lock()
	defer mu.Unlock()
	return rerendering
}

// SetRerendering sets the rerendering flag.
func SetRerendering(v bool) {
	mu.Lock()
	defer mu.Unlock()
	rerendering = v
}

// CaptureRerenderable stores a Rerenderable (e.g., a *table.Table) and visible columns
// so the output can be re-rendered with fresh data after --wait completes.
func CaptureRerenderable(r Rerenderable, visibleCols []string) {
	mu.Lock()
	defer mu.Unlock()
	lastRerenderable = r
	lastVisibleCols = visibleCols
}

// GetRerenderable returns the captured Rerenderable and its visible columns, or nil if not set.
func GetRerenderable() (Rerenderable, []string) {
	mu.Lock()
	defer mu.Unlock()
	return lastRerenderable, lastVisibleCols
}

// Reset clears all stored state. Call this before each command execution if needed.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	lastHref = ""
	lastRenderInfo = nil
	lastRerenderable = nil
	lastVisibleCols = nil
	rerendering = false
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

// FetchResource performs a GET request on the captured href and returns the
// parsed JSON response. Used to re-fetch a resource after waiting so we can
// re-render the output with the final state.
func FetchResource() (any, error) {
	href := GetHref()
	if href == "" {
		return nil, fmt.Errorf("no href captured")
	}

	fullURL := buildFullURL(href)

	cl, err := client.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	cfg := cl.CloudClient.GetConfig()

	return fetchJSON(fullURL, cfg.Token, cfg.Username, cfg.Password)
}

// fetchJSON performs a GET request and returns the parsed JSON response.
func fetchJSON(url, token, username, password string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else if username != "" {
		req.SetBasicAuth(username, password)
	}
	userAgent := viper.GetString(constants.CLIHttpUserAgent)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	httpClient := &http.Client{Timeout: httpTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
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
