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
	neturl "net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"
)

var pollInterval = 5 * time.Second

const httpTimeout = 10 * time.Second

const (
	requestProgressTpl = `{{ etime . }} {{ "Waiting for request" }}{{ cycle . "." ".." "..." "...."}}`
	stateProgressTpl   = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".." "..." "...."}}`
)

// Rerenderable can re-render its output with fresh source data.
// Implemented by *table.Table without requiring an import of that package.
type Rerenderable interface {
	Extract(sourceData any) error
	Render(visibleCols []string) (string, error)
}

var (
	mu               sync.Mutex
	lastHref         string
	lastMethod       string // HTTP method of the captured request (POST, DELETE, etc.)
	lastRequestURL   string // Location header from response (request status URL)
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
func CaptureRequestURL(method, url, locationHeader string) {
	mu.Lock()
	defer mu.Unlock()
	lastMethod = method
	// Always capture request status URL (Location header) so we can poll
	// the request to completion before polling resource state.
	if locationHeader != "" {
		lastRequestURL = locationHeader
	}
	// Only capture resource URL if no href was already set from table output.
	// Table-captured hrefs are more accurate since they come from the response body.
	if lastHref == "" {
		lastHref = url
	}
}

// IsDeleteOperation returns true if the captured HTTP method was DELETE.
func IsDeleteOperation() bool {
	mu.Lock()
	defer mu.Unlock()
	return lastMethod == http.MethodDelete
}

// GetRequestStatusURL returns the captured request status URL (from Location header).
func GetRequestStatusURL() string {
	mu.Lock()
	defer mu.Unlock()
	return lastRequestURL
}

// Reset clears all stored state.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	lastHref = ""
	lastMethod = ""
	lastRequestURL = ""
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
	m := extractMap(sourceData)
	if m == nil {
		return ""
	}
	if href, ok := m["href"].(string); ok {
		return href
	}
	return ""
}

// ExtractID extracts the top-level "id" field from sourceData.
// Used to build resource URLs for APIs that don't include href in responses
// (e.g. postgres-v1, mongo).
func ExtractID(sourceData any) string {
	m := extractMap(sourceData)
	if m == nil {
		return ""
	}
	if id, ok := m["id"].(string); ok {
		return id
	}
	return ""
}

func extractMap(sourceData any) map[string]any {
	b, err := json.Marshal(sourceData)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	// Skip list/collection responses
	if _, hasItems := m["items"]; hasItems {
		return nil
	}
	return m
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

	// If we have a request status URL (from Location header), poll it first.
	// This ensures the API request completes before we check resource state,
	// which prevents reading stale data after updates and avoids polling
	// action endpoints (start/stop/reboot) that don't support GET.
	if reqURL := GetRequestStatusURL(); reqURL != "" {
		if err := pollURL(ctx, w, requestProgressTpl, reqURL, token, username, password, false); err != nil {
			return err
		}
	}

	// For action commands (start/stop/reboot/suspend/resume), the captured
	// href is the action endpoint (e.g. /servers/{id}/start) which doesn't
	// support GET. After polling the request status above, we're done.
	if isActionEndpoint(href) {
		return nil
	}

	// Collect all URLs to poll: the resource itself + all parent resources
	urls := resourceAndParentURLs(href)
	isDelete := IsDeleteOperation()

	for i, url := range urls {
		fullURL := buildFullURL(url)
		// Only the first URL (the resource itself) might be a delete.
		// Parent resources are never being deleted, so 404 on them is always transient.
		deleteOp := isDelete && i == 0

		if err := pollURL(ctx, w, stateProgressTpl, fullURL, token, username, password, deleteOp); err != nil {
			return err
		}
	}

	return nil
}

// pollURL polls a single URL with progress bar (text mode) or silently (JSON mode).
func pollURL(ctx context.Context, w io.Writer, tpl, url, token, username, password string, isDelete bool) error {
	if isStructuredOutput() {
		return Poll(ctx, url, token, username, password, isDelete)
	}

	bar := pb.New(1)
	bar.SetWriter(w)
	bar.SetTemplateString(tpl)
	bar.Start()

	err := Poll(ctx, url, token, username, password, isDelete)
	if err != nil {
		bar.SetTemplateString(tpl + " FAILED")
		bar.Finish()
		return err
	}
	bar.SetTemplateString(tpl + " DONE")
	bar.Finish()
	return nil
}

// isActionEndpoint returns true for server action endpoints that don't support GET.
func isActionEndpoint(href string) bool {
	parts := strings.Split(strings.TrimRight(href, "/"), "/")
	if len(parts) == 0 {
		return false
	}
	last := parts[len(parts)-1]
	switch last {
	case "start", "stop", "reboot", "suspend", "resume":
		return true
	}
	return false
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
	hc.Transport = &capturingTransport{wrapped: transport, waitEnabled: viper.GetBool(constants.ArgWait)}
}

// capturingTransport wraps an http.RoundTripper and captures the request URL
// from mutating HTTP methods (POST, PUT, PATCH, DELETE) into globalwait state.
type capturingTransport struct {
	wrapped     http.RoundTripper
	waitEnabled bool
}

func (t *capturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.wrapped.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Only capture URLs from mutating methods
	switch req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		if t.waitEnabled {
			CaptureRequestURL(req.Method, req.URL.String(), resp.Header.Get("Location"))
		}
	}

	return resp, err
}

// uuidRegex matches standard UUID format used by IONOS APIs.
var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// looksLikeResourceID returns true if the string looks like a resource ID
// (UUID or pure numeric). Uses strict UUID regex to avoid false positives
// on hyphenated resource type names like "private-cross-connects".
func looksLikeResourceID(s string) bool {
	if s == "" {
		return false
	}
	if uuidRegex.MatchString(s) {
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
func Poll(ctx context.Context, url, token, username, password string, isDelete bool) error {
	httpClient := &http.Client{Timeout: httpTimeout}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	userAgent := viper.GetString(constants.CLIHttpUserAgent)
	firstSuccessfulPoll := true

	for {
		// Check state immediately (first iteration), then on each tick
		state, err := fetchState(ctx, httpClient, url, token, username, password, userAgent, isDelete)
		if err != nil {
			// Auth errors are not transient, fail immediately
			if strings.Contains(err.Error(), "authentication failed") {
				return err
			}
			// Other errors (network, 5xx, bad JSON) are transient, retry
		} else if state != "" {
			switch strings.ToUpper(state) {
			case "AVAILABLE", "ACTIVE", "READY", "DONE":
				return nil
			case "FAILED":
				return fmt.Errorf("resource at %s entered FAILED state", url)
			}
			firstSuccessfulPoll = false
		} else if firstSuccessfulPoll {
			// First successful poll returned no state field. Resource
			// doesn't track provisioning state (e.g. /um/groups). Treat
			// as immediately ready.
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for %s to become AVAILABLE", url)
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

	if !strings.HasPrefix(href, "/") {
		href = "/" + href
	}

	return appendDepthParam(strings.TrimRight(baseURL, "/") + href)
}

func appendDepthParam(rawURL string) string {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	q.Set("depth", "1")
	u.RawQuery = q.Encode()
	return u.String()
}

type apiResponse struct {
	Metadata *apiMetadata `json:"metadata"`
}

type apiMetadata struct {
	State  string `json:"state"`
	Status string `json:"status"` // VPN uses "status" instead of "state"
}

func fetchState(ctx context.Context, httpClient *http.Client, url, token, username, password, userAgent string, isDelete bool) (string, error) {
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
				time.Sleep(time.Duration(d) * time.Second)
			}
		}
		return "", fmt.Errorf("rate limited (HTTP 429)")
	}

	// Server errors are transient, will be retried
	if resp.StatusCode >= 500 {
		return "", fmt.Errorf("server error (HTTP %d)", resp.StatusCode)
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

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to fetch resource (HTTP %d)", resp.StatusCode)
	}

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

