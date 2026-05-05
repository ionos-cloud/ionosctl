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

// CaptureRequestURL stores the URL from an API request (e.g. resp.RequestURL).
// When no href was captured from table output (delete/detach commands), this URL
// is used to derive the target resource for --wait polling.
func CaptureRequestURL(url string) {
	mu.Lock()
	defer mu.Unlock()
	// Only capture if no href was already set from table output (create/update).
	// Table-captured hrefs are more accurate since they come from the response body.
	if lastHref == "" {
		lastHref = url
	}
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

// SetResourceHref constructs and captures an href from API path segments.
// Use in commands that don't produce table output (delete, detach) but where
// --wait should poll a parent resource. Example:
//
//	globalwait.SetResourceHref("cloudapi", "v6", "datacenters", dcId, "servers", serverId)
func SetResourceHref(pathSegments ...string) {
	baseURL := viper.GetString(constants.ArgServerUrl)
	if baseURL == "" {
		baseURL = constants.DefaultApiURL
	}
	href := strings.TrimRight(baseURL, "/") + "/" + strings.Join(pathSegments, "/")
	CaptureHref(href)
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
// It then walks up the resource hierarchy and polls each parent until AVAILABLE too.
// Progress output is written to w (typically os.Stderr).
// Returns nil if no href was captured (command doesn't deal with API resources).
func WaitForAvailable(w io.Writer, token, username, password string) error {
	href := GetHref()
	if href == "" {
		return nil
	}

	timeout := time.Duration(viper.GetInt(constants.ArgTimeout)) * time.Second
	if timeout <= 0 {
		timeout = time.Duration(constants.DefaultTimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Collect all URLs to poll: the resource itself + all parent resources
	urls := resourceAndParentURLs(href)

	for _, url := range urls {
		fullURL := buildFullURL(url)

		if isStructuredOutput() {
			if err := pollWithJSONLog(ctx, w, fullURL, token, username, password); err != nil {
				return err
			}
			continue
		}

		bar := pb.New(1)
		bar.SetWriter(w)
		bar.SetTemplateString(progressTpl)
		bar.Start()

		err := Poll(ctx, fullURL, token, username, password)
		if err != nil {
			bar.SetTemplateString(progressTpl + " FAILED")
			bar.Finish()
			return err
		}
		bar.SetTemplateString(progressTpl + " DONE")
		bar.Finish()
	}

	return nil
}

// resourceAndParentURLs returns the given href plus all parent resource hrefs,
// from deepest to shallowest. For example, given:
//
//	https://api.ionos.com/cloudapi/v6/datacenters/dc1/servers/srv1/volumes/vol1
//
// it returns:
//
//	[".../volumes/vol1", ".../servers/srv1", ".../datacenters/dc1"]
//
// Non-CloudAPI hrefs (no /cloudapi/ path) return just the original href.
func resourceAndParentURLs(href string) []string {
	urls := []string{href}

	// Walk up by stripping last two path segments (resource-type/resource-id)
	current := href
	for {
		parent := parentHref(current)
		if parent == "" {
			break
		}
		urls = append(urls, parent)
		current = parent
	}

	return urls
}

// parentHref strips the last two path segments (resource-type/id) to get the
// parent resource href. Returns "" if there's no valid parent.
//
// Works with both CloudAPI and regional API URL structures:
//   - CloudAPI: https://api.ionos.com/cloudapi/v6/datacenters/dc1/servers/srv1
//   - Regional: https://vpn.de-fra.ionos.com/wireguardgateways/gw1/peers/p1
//
// Stops when stripping would leave no resource pair (type+id) after the host.
func parentHref(href string) string {
	parts := strings.Split(href, "/")

	// After splitting, first 3 parts are always: "https:", "", "host"
	// Then optional API prefix segments (e.g. "cloudapi", "v6") followed by
	// resource pairs (type/id). We need at least 2 resource pairs (4 path
	// segments after host) to have a parent.
	//
	// Find where resource path starts by skipping non-UUID/non-resource segments
	// after the host. Simpler: the candidate after stripping 2 must still have
	// at least one type/id pair (2 segments) after the host portion.
	//
	// Minimum: "https:"/""/host/type/id/type/id = 7 parts
	// After strip: "https:"/""/host/type/id = 5 parts (valid parent)
	// Next strip would give: "https:"/""/host = 3 parts (just host, not valid)
	if len(parts) < 7 {
		return ""
	}

	candidate := strings.Join(parts[:len(parts)-2], "/")

	// Candidate must end with a resource ID (last segment should look like
	// an ID, not an API path component like "v6" or "cloudapi").
	// Resource IDs are typically UUIDs or alphanumeric strings.
	lastSeg := parts[len(parts)-3] // last segment of candidate
	if !looksLikeResourceID(lastSeg) {
		return ""
	}

	return candidate
}

// WrapTransport wraps an http.Client's Transport so that every response URL
// is captured for --wait polling. This makes delete/detach commands work
// across all SDK clients without per-command changes.
func WrapTransport(hc *http.Client) {
	if hc == nil {
		return
	}
	if _, ok := hc.Transport.(*capturingTransport); ok {
		return // already wrapped
	}
	transport := hc.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	hc.Transport = &capturingTransport{wrapped: transport}
}

// capturingTransport wraps an http.RoundTripper and captures the request URL
// from mutating HTTP methods (POST, PUT, PATCH, DELETE) into globalwait state.
type capturingTransport struct {
	wrapped http.RoundTripper
}

func (t *capturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.wrapped.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Only capture URLs from mutating methods
	switch req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if viper.GetBool(constants.ArgWait) {
			CaptureRequestURL(req.URL.String())
		}
	}

	return resp, err
}

// looksLikeResourceID returns true if the string looks like a resource ID
// (UUID, numeric, or other alphanumeric ID) rather than an API path component
// (like "cloudapi", "v6", "datacenters").
func looksLikeResourceID(s string) bool {
	if s == "" {
		return false
	}
	// UUIDs contain hyphens and hex chars
	if strings.Contains(s, "-") && len(s) > 8 {
		return true
	}
	// Pure numeric IDs
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// FetchResource performs a GET on the captured href and returns parsed JSON.
// Used to re-fetch a resource after waiting so we can re-render with final state.
func FetchResource(token, username, password string) (any, error) {
	href := GetHref()
	if href == "" {
		return nil, fmt.Errorf("no href captured")
	}

	fullURL := buildFullURL(href)
	return fetchJSON(fullURL, token, username, password)
}

// Poll polls the given URL until the resource reaches a terminal ready state
// (AVAILABLE, ACTIVE, READY, DONE) or a failure state (FAILED).
func Poll(ctx context.Context, url, token, username, password string) error {
	httpClient := &http.Client{Timeout: httpTimeout}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	userAgent := viper.GetString(constants.CLIHttpUserAgent)
	firstPoll := true

	for {
		// Check state immediately (first iteration), then on each tick
		state, err := fetchState(ctx, httpClient, url, token, username, password, userAgent)
		if err == nil {
			if state != "" {
				switch strings.ToUpper(state) {
				case "AVAILABLE", "ACTIVE", "READY", "DONE":
					return nil
				case "FAILED":
					return fmt.Errorf("resource entered FAILED state")
				}
			} else if firstPoll {
				// First poll returned no state field at all. API does not
				// support metadata.state/status (e.g. MongoDB, MariaDB).
				return fmt.Errorf("--wait is not supported for this resource: API response has no metadata state field")
			}
		}
		firstPoll = false

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for resource to become AVAILABLE")
		case <-ticker.C:
		}
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

	// 404 means resource was deleted — treat as terminal success
	if resp.StatusCode == http.StatusNotFound {
		return "DONE", nil
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
