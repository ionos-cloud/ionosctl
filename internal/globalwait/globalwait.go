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
	"os"
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

const progressTpl = `{{ etime . }} {{ "Waiting" }}{{ cycle . "." ".." "..." "...."}}`

// Rerenderable can re-render its output with fresh source data.
// Implemented by *table.Table without requiring an import of that package.
type Rerenderable interface {
	Extract(sourceData any) error
	Render(visibleCols []string) (string, error)
}

// AuthCreds holds authentication credentials for polling requests.
type AuthCreds struct {
	Token, Username, Password string
}

var (
	mu               sync.Mutex
	lastHref         string
	lastMethod       string // HTTP method of the captured request (POST, DELETE, etc.)
	lastRequestURL   string // Location header from response (request status URL)
	lastRerenderable Rerenderable
	lastVisibleCols  []string
	rerendering      bool
	sdkTransport     http.RoundTripper // captured from first WrapTransport call, reused by poller
	captureCount     int               // number of captureRequestURL calls (detects bulk operations)
	lastHrefFromGet  bool              // true when lastHref was set by a GET (lower priority than mutating methods)
)

// captureHref stores the given href for later polling.
// Clears lastHrefFromGet since this is an explicit capture from response body.
func captureHref(href string) {
	mu.Lock()
	defer mu.Unlock()
	lastHref = href
	lastHrefFromGet = false
}

// getHref returns the last captured href.
func getHref() string {
	mu.Lock()
	defer mu.Unlock()
	return lastHref
}

// captureRerenderable stores a Rerenderable (e.g. *table.Table) and visible columns
// so output can be re-rendered with fresh data after --wait completes.
func captureRerenderable(r Rerenderable, visibleCols []string) {
	mu.Lock()
	defer mu.Unlock()
	lastRerenderable = r
	lastVisibleCols = visibleCols
}

// getRerenderable returns the captured Rerenderable and its visible columns, or nil if not set.
func getRerenderable() (Rerenderable, []string) {
	mu.Lock()
	defer mu.Unlock()
	return lastRerenderable, lastVisibleCols
}

// isRerendering returns true during the re-render pass after waiting.
// Used by the BeforeRender hook to allow output through on the second call.
func isRerendering() bool {
	mu.Lock()
	defer mu.Unlock()
	return rerendering
}

// setRerendering sets the rerendering flag.
func setRerendering(v bool) {
	mu.Lock()
	defer mu.Unlock()
	rerendering = v
}

// captureGetURL stores the URL from a GET request for --wait polling.
// GET captures have lower priority than mutating methods: lastHref is only set
// if empty, and lastMethod is only set if no mutating method was captured.
// This prevents PreCmdRun GET lookups (completers, validators) from overriding
// state that a subsequent POST/DELETE should control.
func captureGetURL(url string) {
	mu.Lock()
	defer mu.Unlock()
	if lastHref == "" {
		lastHref = url
		lastHrefFromGet = true
	}
	if lastMethod == "" {
		lastMethod = http.MethodGet
	}
}

// captureRequestURL stores the URL from an API request (e.g. resp.RequestURL).
// When no href was captured from table output (delete/detach commands), this URL
// is used to derive the target resource for --wait polling.
func captureRequestURL(method, url, locationHeader string) {
	mu.Lock()
	defer mu.Unlock()
	captureCount++
	lastMethod = method
	// Always capture request status URL (Location header) so we can poll
	// the request to completion before polling resource state.
	if locationHeader != "" {
		lastRequestURL = locationHeader
	}
	// Capture resource URL if no href was already set from table output,
	// or if the current href was set by a GET (lower priority).
	// Table-captured hrefs are more accurate since they come from the response body.
	if lastHref == "" || lastHrefFromGet {
		lastHref = url
		lastHrefFromGet = false
	}
}

// isDeleteOperation returns true if the captured HTTP method was DELETE.
func isDeleteOperation() bool {
	mu.Lock()
	defer mu.Unlock()
	return lastMethod == http.MethodDelete
}

// isGetOperation returns true if the captured HTTP method was GET.
func isGetOperation() bool {
	mu.Lock()
	defer mu.Unlock()
	return lastMethod == http.MethodGet
}

// getRequestStatusURL returns the captured request status URL (from Location header).
func getRequestStatusURL() string {
	mu.Lock()
	defer mu.Unlock()
	return lastRequestURL
}

// getCaptureCount returns how many times captureRequestURL was called since last Reset.
// Used to detect bulk operations (e.g. --all delete) where only the last resource is polled.
func getCaptureCount() int {
	mu.Lock()
	defer mu.Unlock()
	return captureCount
}

// Reset clears all captured state. Call between multiple mutating API calls
// within a single command to prevent mismatched state (e.g. server create
// followed by --promote-volume). Each call to captureRequestURL overwrites
// the request status URL but preserves the first resource href, so without
// Reset() the poller may poll a request status URL from call #2 while using
// a resource href from call #1.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	lastHref = ""
	lastMethod = ""
	lastRequestURL = ""
	lastRerenderable = nil
	lastVisibleCols = nil
	rerendering = false
	sdkTransport = nil
	captureCount = 0
	lastHrefFromGet = false
}

// extractHref extracts the top-level "href" field from sourceData.
// Returns empty string if sourceData is a list (has "items" key), has no href,
// or cannot be parsed.
func extractHref(sourceData any) string {
	m := extractMap(sourceData)
	if m == nil {
		return ""
	}
	if href, ok := m["href"].(string); ok {
		return href
	}
	return ""
}

// extractID extracts the top-level "id" field from sourceData.
// Used to build resource URLs for APIs that don't include href in responses
// (e.g. postgres-v1, mongo).
func extractID(sourceData any) string {
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

// HandleBeforeRender processes table output for --wait. Returns true to render
// normally, false to suppress (output will be re-rendered after wait completes).
// Called from the table.BeforeRender hook adapter in commands/root.go.
func HandleBeforeRender(sourceData any, visibleCols []string, r Rerenderable) bool {
	if !viper.GetBool(constants.ArgWait) || isRerendering() {
		return true // render normally
	}
	// Only suppress output for known valid formats. Invalid formats
	// (e.g. typo "-o jso") should render normally so the error surfaces
	// immediately instead of being lost after wait + re-render failure.
	switch viper.GetString(constants.ArgOutput) {
	case "text", "json", "api-json":
	default:
		return true
	}
	href := extractHref(sourceData)
	if href == "" {
		// No href in response (e.g. postgres-v1, mongo, DNS).
		id := extractID(sourceData)
		if id == "" {
			return true // list or unrecognized format - render normally
		}
		// For GET, the transport-captured URL is already the resource URL.
		// For POST/PUT/PATCH, it's the collection URL - append the id.
		if base := getHref(); base != "" && !isGetOperation() {
			captureHref(strings.TrimRight(base, "/") + "/" + id)
		}
		if getHref() == "" {
			return true // no href and no fallback, render normally
		}
	} else {
		// Response has href, use it directly. More specific than the
		// transport-captured URL. buildFullURL resolves relative hrefs.
		captureHref(href)
	}
	captureRerenderable(r, visibleCols)
	return false // suppress initial output
}

// WaitAndRerender polls until the resource is available, then re-renders output
// with fresh data showing the final state. Call after successful command execution
// when --wait is set. Progress and warnings are written to stderr; re-rendered
// output is written to stdout.
func WaitAndRerender(stderr, stdout io.Writer, creds AuthCreds, quiet bool) error {
	if err := WaitForAvailable(stderr, creds.Token, creds.Username, creds.Password); err != nil {
		return err
	}

	if quiet {
		return nil
	}

	r, cols := getRerenderable()
	if r == nil {
		return nil
	}

	freshData, err := fetchResource(creds.Token, creds.Username, creds.Password)
	if err != nil {
		fmt.Fprintf(stderr, "Warning: could not fetch updated resource: %v\n", err)
		return nil
	}

	setRerendering(true)
	defer setRerendering(false)

	if err := r.Extract(freshData); err != nil {
		fmt.Fprintf(stderr, "Warning: could not extract fresh data: %v\n", err)
		return nil
	}

	out, err := r.Render(cols)
	if err != nil {
		fmt.Fprintf(stderr, "Warning: could not re-render output: %v\n", err)
		return nil
	}

	fmt.Fprint(stdout, out)
	return nil
}

// WaitForAvailable polls the captured href until the resource reaches a terminal ready state.
// It then walks up the resource hierarchy and polls each parent until AVAILABLE too.
// Progress output is written to w (typically os.Stderr).
// Returns nil if no href was captured (command doesn't deal with API resources).
func WaitForAvailable(w io.Writer, token, username, password string) error {
	href := getHref()
	if href == "" {
		return nil
	}

	timeoutSec := viper.GetInt(constants.ArgTimeout)
	if timeoutSec <= 0 {
		if timeoutSec == 0 {
			fmt.Fprintf(w, "Warning: --timeout 0 is not supported, using default %ds\n", constants.DefaultTimeoutSeconds)
		}
		timeoutSec = constants.DefaultTimeoutSeconds
	}
	timeout := time.Duration(timeoutSec) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Collect all URLs to poll in order:
	// 1. Request status URL (Location header) if available
	// 2. Resource URL + parent URLs (unless action endpoint)
	type pollTarget struct {
		url      string
		isDelete bool
	}
	var targets []pollTarget

	if reqURL := getRequestStatusURL(); reqURL != "" {
		targets = append(targets, pollTarget{url: reqURL})
	}

	// Action endpoints (start/stop/reboot/suspend/resume) don't support GET.
	// Only the request status poll above is needed.
	if !isActionEndpoint(href) {
		urls := resourceAndParentURLs(href)
		isDelete := isDeleteOperation()
		for i, url := range urls {
			targets = append(targets, pollTarget{
				url:      buildFullURL(url),
				isDelete: isDelete && i == 0,
			})
		}
	}

	if len(targets) == 0 {
		fmt.Fprintf(w, "Warning: --wait active but no resource URL could be determined for polling\n")
		return nil
	}

	if n := getCaptureCount(); n > 1 {
		fmt.Fprintf(w, "Warning: --wait only polls the last resource from %d operations. For guaranteed completion, run operations individually with --wait.\n", n)
	}

	p := newPoller(token, username, password)

	// Single progress bar for all polls
	if !isStructuredOutput() {
		bar := pb.New(1)
		bar.SetWriter(w)
		bar.SetTemplateString(progressTpl)
		bar.Start()
		defer func() {
			bar.Finish()
			fmt.Fprintln(w)
		}()

		for _, t := range targets {
			if err := p.poll(ctx, t.url, t.isDelete); err != nil {
				bar.SetTemplateString(progressTpl + " FAILED")
				return err
			}
		}
		bar.SetTemplateString(progressTpl + " DONE")
		return nil
	}

	// JSON mode: poll silently
	for _, t := range targets {
		if err := p.poll(ctx, t.url, t.isDelete); err != nil {
			return err
		}
	}
	return nil
}

// isActionEndpoint returns true for action endpoints that don't support GET
// (e.g. server start/stop, database restore, DNS zone transfer).
func isActionEndpoint(href string) bool {
	u, err := neturl.Parse(href)
	if err != nil {
		return false
	}
	parts := strings.Split(strings.TrimRight(u.Path, "/"), "/")
	if len(parts) == 0 {
		return false
	}
	last := parts[len(parts)-1]
	switch last {
	case "start", "stop", "reboot", "suspend", "resume", "restore", "transfer":
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
	mu.Lock()
	if sdkTransport == nil {
		sdkTransport = transport // reuse TLS config in poller
	}
	mu.Unlock()
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

	// Capture URLs when --wait is active.
	// Read viper at call time (not cached) so deprecated flag mapping works.
	if viper.GetBool(constants.ArgWait) {
		switch req.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			captureRequestURL(req.Method, req.URL.String(), resp.Header.Get("Location"))
		case http.MethodGet:
			captureGetURL(req.URL.String())
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

// fetchResource performs a GET on the captured href and returns parsed JSON.
// Used to re-fetch a resource after waiting so we can re-render with final state.
func fetchResource(token, username, password string) (any, error) {
	href := getHref()
	if href == "" {
		return nil, fmt.Errorf("no href captured")
	}

	p := newPoller(token, username, password)
	return p.fetchJSON(buildFullURL(href))
}

// pollURL polls the given URL until the resource reaches a terminal ready state
// (AVAILABLE, ACTIVE, READY, DONE) or a failure state (FAILED).
func pollURL(ctx context.Context, url, token, username, password string, isDelete bool) error {
	return newPoller(token, username, password).poll(ctx, url, isDelete)
}

// poller holds shared HTTP client, auth, and user-agent for polling requests.
type poller struct {
	client    *http.Client
	token     string
	username  string
	password  string
	userAgent string
}

func newPoller(token, username, password string) *poller {
	mu.Lock()
	transport := sdkTransport
	mu.Unlock()
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &poller{
		client:    &http.Client{Timeout: httpTimeout, Transport: transport},
		token:     token,
		username:  username,
		password:  password,
		userAgent: viper.GetString(constants.CLIHttpUserAgent),
	}
}

func (p *poller) poll(ctx context.Context, url string, isDelete bool) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	firstSuccessfulPoll := true

	for {
		state, err := p.fetchState(ctx, url, isDelete)
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
				return fmt.Errorf("resource at %s entered %s state", url, state)
			case "DESTROYING":
				if !isDelete {
					return fmt.Errorf("resource at %s entered %s state", url, state)
				}
				// Delete in progress - keep polling until 404 returns "DONE"
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
	depth := viper.GetInt(constants.FlagDepth)
	if depth <= 0 {
		depth = 1
	}
	q := u.Query()
	q.Set("depth", strconv.Itoa(depth))
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

func isStructuredOutput() bool {
	switch viper.GetString(constants.ArgOutput) {
	case "json", "api-json":
		return true
	default:
		return false
	}
}
