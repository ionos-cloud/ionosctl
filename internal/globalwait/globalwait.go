// Package globalwait provides a global --wait mechanism for ionosctl.
// When --wait is set, it captures the href from the command's API response,
// then polls that href until the resource reaches a terminal ready state.
//
// This package intentionally has no dependency on the table package.
// The Rerenderable interface is satisfied by *table.Table implicitly,
// and wiring is done in commands/root.go via the table.BeforeRender hook.
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

const httpTimeout = 10 * time.Second

const progressTpl = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".." "..." "...."}}`

// Rerenderable can re-render its output with fresh source data.
// Implemented by *table.Table without requiring an import of that package.
type Rerenderable interface {
	Extract(sourceData any) error
	Render(visibleCols []string) (string, error)
}

var (
	mu               sync.Mutex
	lastHref         string
	lastRerenderable Rerenderable
	lastVisibleCols  []string
	rerendering      bool
)

// CaptureHref stores the given href for later polling.
func CaptureHref(href string) {
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

// CaptureRerenderable stores a Rerenderable (e.g. *table.Table) and visible columns
// so output can be re-rendered with fresh data after --wait completes.
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

// IsRerendering returns true during the re-render pass after waiting.
// Used by the BeforeRender hook to allow output through on the second call.
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

// Reset clears all stored state.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	lastHref = ""
	lastRerenderable = nil
	lastVisibleCols = nil
	rerendering = false
}

// ExtractHref extracts the top-level "href" field from sourceData.
// Returns empty string if sourceData is a list (has "items" key), has no href,
// or cannot be parsed.
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

// WaitForAvailable polls the captured href until the resource reaches a terminal ready state.
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

	fullURL := buildFullURL(href)

	cl, err := client.Get()
	if err != nil {
		return fmt.Errorf("failed to get client for wait polling: %w", err)
	}
	cfg := cl.CloudClient.GetConfig()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if isStructuredOutput() {
		return pollWithJSONLog(ctx, w, fullURL, cfg.Token, cfg.Username, cfg.Password)
	}

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

// FetchResource performs a GET on the captured href and returns parsed JSON.
// Used to re-fetch a resource after waiting so we can re-render with final state.
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
		}
		// BUSY, DEPLOYING, UPDATING, PROVISIONING, DESTROYING etc. - keep polling
	}
}

// --- Internal helpers ---

func buildFullURL(href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return appendDepthParam(href)
	}

	baseURL := viper.GetString(constants.ArgServerUrl)
	if baseURL == "" {
		baseURL = constants.DefaultApiURL
	}

	return appendDepthParam(strings.TrimRight(baseURL, "/") + href)
}

func appendDepthParam(url string) string {
	if strings.Contains(url, "?") {
		return url + "&depth=1"
	}
	return url + "?depth=1"
}

type apiResponse struct {
	Metadata *apiMetadata `json:"metadata"`
}

type apiMetadata struct {
	State  string `json:"state"`
	Status string `json:"status"` // VPN uses "status" instead of "state"
}

func fetchState(ctx context.Context, httpClient *http.Client, url, token, username, password, userAgent string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	setAuth(req, token, username, password)
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

	state := body.Metadata.State
	if state == "" {
		state = body.Metadata.Status
	}
	return state, nil
}

func fetchJSON(url, token, username, password string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	setAuth(req, token, username, password)

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

func setAuth(req *http.Request, token, username, password string) {
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else if username != "" {
		req.SetBasicAuth(username, password)
	}
}

func isStructuredOutput() bool {
	switch viper.GetString(constants.ArgOutput) {
	case "json", "api-json":
		return true
	default:
		return false
	}
}

func pollWithJSONLog(ctx context.Context, w io.Writer, url, token, username, password string) error {
	logJSON(w, "Waiting for state...")
	err := Poll(ctx, url, token, username, password)
	if err != nil {
		logJSON(w, "FAILED")
		return err
	}
	logJSON(w, "DONE")
	return nil
}

func logJSON(w io.Writer, msg string) {
	out, _ := json.Marshal(msg)
	fmt.Fprintln(w, string(out))
}
