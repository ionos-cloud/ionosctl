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
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"
)

var pollInterval = 5 * time.Second

const httpTimeout = 10 * time.Second

// ProgressTpl is the progress bar template used for --wait polling.
const ProgressTpl = `{{ etime . }} {{ "Waiting" }}{{ cycle . "." ".." "..." "...."}}`

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

// Waiter holds all captured state for a single --wait lifecycle.
// Use defaultWaiter (via package-level functions) for normal CLI operation,
// or create a local instance for isolated testing.
type Waiter struct {
	mu            sync.Mutex
	href          string
	method        string // HTTP method of the captured request (POST, DELETE, etc.)
	requestURL    string // Location header from response (request status URL)
	rerenderable  Rerenderable
	visibleCols   []string
	rerendering   bool
	transport     http.RoundTripper // captured from first WrapTransport call, reused by poller
	captureCount  int               // number of captureRequestURL calls (detects bulk operations)
	hrefFromGet   bool              // true when href was set by a GET (lower priority than mutating methods)
	initialOutput string            // buffered initial output for fallback when re-render fails
	done          bool              // set by MarkDone to skip post-command WaitAndRerender
}

var defaultWaiter = &Waiter{}

// --- Waiter methods: state getters/setters ---

func (w *Waiter) captureHref(href string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.href = href
	w.hrefFromGet = false
}

func (w *Waiter) getHref() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.href
}

func (w *Waiter) captureRerenderable(r Rerenderable, visibleCols []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rerenderable = r
	w.visibleCols = visibleCols
}

func (w *Waiter) getRerenderable() (Rerenderable, []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.rerenderable, w.visibleCols
}

func (w *Waiter) getInitialOutput() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.initialOutput
}

func (w *Waiter) isRerendering() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.rerendering
}

func (w *Waiter) setRerendering(v bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rerendering = v
}

// captureGetURL stores the URL from a GET request for --wait polling.
// GET captures have lower priority than mutating methods: href is only set
// if empty, and method is only set if no mutating method was captured.
// This prevents PreCmdRun GET lookups (completers, validators) from overriding
// state that a subsequent POST/DELETE should control.
func (w *Waiter) captureGetURL(url string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.href == "" {
		if u, err := neturl.Parse(url); err == nil {
			u.RawQuery = ""
			w.href = u.String()
		} else {
			w.href = url
		}
		w.hrefFromGet = true
	}
	if w.method == "" {
		w.method = http.MethodGet
	}
}

// captureRequestURL stores the URL from an API request (e.g. resp.RequestURL).
// When no href was captured from table output (delete/detach commands), this URL
// is used to derive the target resource for --wait polling.
func (w *Waiter) captureRequestURL(method, url, locationHeader string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.captureCount++
	w.method = method
	// Always capture request status URL (Location header) so we can poll
	// the request to completion before polling resource state.
	if locationHeader != "" {
		w.requestURL = locationHeader
	}
	// Capture resource URL if no href was already set from table output,
	// or if the current href was set by a GET (lower priority).
	// Table-captured hrefs are more accurate since they come from the response body.
	// Strip query parameters: SDK clients add ?depth=&limit=&offset= to request
	// URLs, and these are invalid when used for polling (cause HTTP 400).
	if w.href == "" || w.hrefFromGet {
		if u, err := neturl.Parse(url); err == nil {
			u.RawQuery = ""
			w.href = u.String()
		} else {
			w.href = url
		}
		w.hrefFromGet = false
	}
}

func (w *Waiter) isPostOperation() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.method == http.MethodPost
}

func (w *Waiter) isDeleteOperation() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.method == http.MethodDelete
}

func (w *Waiter) isGetOperation() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.method == http.MethodGet
}

func (w *Waiter) getRequestStatusURL() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.requestURL
}

func (w *Waiter) getCaptureCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.captureCount
}

func (w *Waiter) isDone() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.done
}

// --- Waiter methods: exported ---

// Reset clears all captured state. Call between multiple mutating API calls
// within a single command to prevent mismatched state (e.g. server create
// followed by --promote-volume). Each call to captureRequestURL overwrites
// the request status URL but preserves the first resource href, so without
// Reset() the poller may poll a request status URL from call #2 while using
// a resource href from call #1.
func (w *Waiter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.href = ""
	w.method = ""
	w.requestURL = ""
	w.rerenderable = nil
	w.visibleCols = nil
	w.rerendering = false
	w.transport = nil
	w.captureCount = 0
	w.hrefFromGet = false
	w.initialOutput = ""
}

// MarkDone signals that all waiting has been handled inline by the command.
// HandleBeforeRender will render normally, and WaitAndRerender will be a no-op.
// Used by commands that perform multiple mutating API calls and manage their
// own wait lifecycle (e.g. server create --promote-volume).
func (w *Waiter) MarkDone() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.done = true
}

// HandleBeforeRender processes table output for --wait. Returns true to render
// normally, false to suppress (output will be re-rendered after wait completes).
// Called from the table.BeforeRender hook adapter in commands/root.go.
func (w *Waiter) HandleBeforeRender(sourceData any, visibleCols []string, r Rerenderable) bool {
	if !viper.GetBool(constants.ArgWait) || w.isRerendering() || w.isDone() {
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
		// For GET/PUT/PATCH, the transport-captured URL is already the resource URL.
		// For POST, it's the collection URL - append the id to form the resource URL.
		if base := w.getHref(); base != "" {
			if w.isPostOperation() {
				w.captureHref(strings.TrimRight(base, "/") + "/" + id)
			}
			// else: PUT/PATCH/GET already have the resource URL, keep as-is
		}
		if w.getHref() == "" {
			return true // no href and no fallback, render normally
		}
	} else {
		// Response has href, use it directly. More specific than the
		// transport-captured URL. buildFullURL resolves relative hrefs.
		w.captureHref(href)
	}
	w.captureRerenderable(r, visibleCols)

	// Pre-render and buffer the initial output so we can fall back to it
	// if re-rendering fails after wait (e.g., fetchResource gets 429/5xx).
	// Setting rerendering=true prevents infinite recursion: the recursive
	// HandleBeforeRender call sees isRerendering()=true and returns true,
	// allowing the inner Render to produce the actual output.
	w.setRerendering(true)
	if initial, err := r.Render(visibleCols); err == nil {
		w.mu.Lock()
		w.initialOutput = initial
		w.mu.Unlock()
	}
	w.setRerendering(false)

	return false // suppress initial output
}

// WaitAndRerender polls until the resource is available, then re-renders output
// with fresh data showing the final state. Call after successful command execution
// when --wait is set. Progress and warnings are written to stderr; re-rendered
// output is written to stdout.
func (w *Waiter) WaitAndRerender(stderr, stdout io.Writer, creds AuthCreds, quiet bool) error {
	if w.isDone() {
		return nil
	}

	r, cols := w.getRerenderable()
	if r == nil && w.getRequestStatusURL() == "" {
		// No output was suppressed (e.g. list command) and no async request
		// to track — nothing to wait for or re-render.
		return nil
	}

	if err := w.WaitForAvailable(stderr, creds.Token, creds.Username, creds.Password); err != nil {
		return err
	}

	if quiet || r == nil {
		return nil
	}

	freshData, err := w.fetchResource(creds.Token, creds.Username, creds.Password)
	if err != nil {
		fmt.Fprintf(stderr, "Warning: could not fetch updated resource: %v\n", err)
		return w.emitFallbackOutput(stderr, stdout)
	}

	w.setRerendering(true)
	defer w.setRerendering(false)

	if err := r.Extract(freshData); err != nil {
		fmt.Fprintf(stderr, "Warning: could not extract fresh data: %v\n", err)
		return w.emitFallbackOutput(stderr, stdout)
	}

	out, err := r.Render(cols)
	if err != nil {
		fmt.Fprintf(stderr, "Warning: could not re-render output: %v\n", err)
		return w.emitFallbackOutput(stderr, stdout)
	}

	fmt.Fprint(stdout, out)
	return nil
}

// emitFallbackOutput writes the buffered initial output to stdout when re-rendering
// fails. This prevents the command from exiting with zero stdout output after
// HandleBeforeRender suppressed the initial render. If no fallback is available,
// returns an error so the caller can exit non-zero.
func (w *Waiter) emitFallbackOutput(stderr, stdout io.Writer) error {
	if initial := w.getInitialOutput(); initial != "" {
		fmt.Fprintf(stderr, "Warning: showing pre-wait output (may not reflect final state)\n")
		_, err := fmt.Fprint(stdout, initial)
		return err
	}
	return fmt.Errorf("--wait re-render failed and no fallback output available")
}

// WaitForAvailable polls the captured href until the resource reaches a terminal ready state.
// It then walks up the resource hierarchy and polls each parent until AVAILABLE too.
// Progress output is written to wr (typically os.Stderr).
// Returns nil if no href was captured (command doesn't deal with API resources).
func (w *Waiter) WaitForAvailable(wr io.Writer, token, username, password string) error {
	href := w.getHref()
	if href == "" {
		return nil
	}

	// Collect all URLs to poll in order:
	// 1. Request status URL (Location header) if available
	// 2. Resource URL + parent URLs (unless action endpoint)
	type pollTarget struct {
		url      string
		isDelete bool
	}
	var targets []pollTarget

	reqURL := w.getRequestStatusURL()

	// Action endpoints (start/stop/reboot/suspend/resume) don't support GET.
	// Only poll the request status URL if available; otherwise nothing to do.
	if isActionEndpoint(href) {
		if reqURL == "" {
			fmt.Fprintf(wr, "Warning: --wait: API response missing Location header for %s, cannot track request status\n", href)
			return nil
		}
		targets = append(targets, pollTarget{url: reqURL})
	} else {
		if reqURL != "" {
			targets = append(targets, pollTarget{url: reqURL})
		}
		urls := resourceAndParentURLs(href)
		isDelete := w.isDeleteOperation()
		for i, url := range urls {
			targets = append(targets, pollTarget{
				url:      buildFullURL(url),
				isDelete: isDelete && i == 0,
			})
		}
	}

	if len(targets) == 0 {
		fmt.Fprintf(wr, "Warning: --wait active but no resource URL could be determined for polling\n")
		return nil
	}

	timeoutSec := viper.GetInt(constants.ArgTimeout)
	if timeoutSec <= 0 {
		fmt.Fprintf(wr, "Warning: --timeout %d is not supported, using default %ds\n", timeoutSec, constants.DefaultTimeoutSeconds)
		timeoutSec = constants.DefaultTimeoutSeconds
	}
	timeout := time.Duration(timeoutSec) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if n := w.getCaptureCount(); n > 1 {
		fmt.Fprintf(wr, "Warning: --wait only polls the last resource from %d operations. For guaranteed completion, run operations individually with --wait.\n", n)
	}

	p := w.newPoller(token, username, password)

	// Single progress bar for all polls
	if !isStructuredOutput() {
		bar := pb.New(1)
		bar.SetWriter(wr)
		bar.SetTemplateString(ProgressTpl)
		bar.Start()
		defer bar.Finish()

		for _, t := range targets {
			if err := p.poll(ctx, t.url, t.isDelete); err != nil {
				var pf *provisioningFailure
				if errors.As(err, &pf) {
					bar.SetTemplateString(ProgressTpl + " " + pf.state)
					fmt.Fprintf(wr, "\nWarning: %v\n", pf)
					return nil
				}
				bar.SetTemplateString(ProgressTpl + " FAILED")
				return err
			}
		}
		bar.SetTemplateString(ProgressTpl + " DONE")
		return nil
	}

	// JSON mode: poll silently
	for _, t := range targets {
		if err := p.poll(ctx, t.url, t.isDelete); err != nil {
			var pf *provisioningFailure
			if errors.As(err, &pf) {
				fmt.Fprintf(wr, "Warning: %v\n", pf)
				return nil
			}
			return err
		}
	}
	return nil
}

// WrapTransport wraps an http.Client's Transport so that every response URL
// is captured for --wait polling. This makes delete/detach commands work
// across all SDK clients without per-command changes.
func (w *Waiter) WrapTransport(hc *http.Client) {
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
	w.mu.Lock()
	if w.transport == nil {
		w.transport = transport // reuse TLS config in poller
	}
	w.mu.Unlock()
	hc.Transport = &capturingTransport{wrapped: transport, waiter: w}
}

// fetchResource performs a GET on the captured href and returns parsed JSON.
// Used to re-fetch a resource after waiting so we can re-render with final state.
func (w *Waiter) fetchResource(token, username, password string) (any, error) {
	href := w.getHref()
	if href == "" {
		return nil, fmt.Errorf("no href captured")
	}

	p := w.newPoller(token, username, password)
	return p.fetchJSON(buildFullURL(href))
}

// pollURL polls the given URL until the resource reaches a terminal ready state
// (AVAILABLE, ACTIVE, READY, DONE) or a failure state (FAILED).
func (w *Waiter) pollURL(ctx context.Context, url, token, username, password string, isDelete bool) error {
	return w.newPoller(token, username, password).poll(ctx, url, isDelete)
}

func (w *Waiter) newPoller(token, username, password string) *poller {
	w.mu.Lock()
	transport := w.transport
	w.mu.Unlock()
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

// --- Package-level delegates (backward-compatible API) ---

// Reset clears all captured state on the default Waiter.
func Reset() { defaultWaiter.Reset() }

// MarkDone signals that all waiting has been handled inline by the command.
func MarkDone() { defaultWaiter.MarkDone() }

// HandleBeforeRender processes table output for --wait on the default Waiter.
func HandleBeforeRender(sourceData any, visibleCols []string, r Rerenderable) bool {
	return defaultWaiter.HandleBeforeRender(sourceData, visibleCols, r)
}

// WaitAndRerender polls and re-renders on the default Waiter.
func WaitAndRerender(stderr, stdout io.Writer, creds AuthCreds, quiet bool) error {
	return defaultWaiter.WaitAndRerender(stderr, stdout, creds, quiet)
}

// WaitForAvailable polls the captured href on the default Waiter.
func WaitForAvailable(w io.Writer, token, username, password string) error {
	return defaultWaiter.WaitForAvailable(w, token, username, password)
}

// WrapTransport wraps an http.Client's Transport on the default Waiter.
func WrapTransport(hc *http.Client) { defaultWaiter.WrapTransport(hc) }

// --- capturingTransport ---

// capturingTransport wraps an http.RoundTripper and captures the request URL
// from mutating HTTP methods (POST, PUT, PATCH, DELETE) into Waiter state.
type capturingTransport struct {
	wrapped http.RoundTripper
	waiter  *Waiter
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
			t.waiter.captureRequestURL(req.Method, req.URL.String(), resp.Header.Get("Location"))
		case http.MethodGet:
			t.waiter.captureGetURL(req.URL.String())
		}
	}

	return resp, err
}
