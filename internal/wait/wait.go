package wait

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
)

// WaitOption is a function that configures WaitOptions.
type WaitOption func(*WaitOptions)

// WaitOptions holds configuration for waiting behavior.
type WaitOptions struct {
	PollInterval time.Duration
	Ctx          context.Context
}

// WithPollInterval returns a WaitOption that sets the poll interval.
func WithPollInterval(d time.Duration) WaitOption {
	return func(opts *WaitOptions) {
		opts.PollInterval = d
	}
}

// WithContext returns a WaitOption that sets the context.
func WithContext(ctx context.Context) WaitOption {
	return func(opts *WaitOptions) {
		opts.Ctx = ctx
	}
}

// For polls the provided href until the resource reaches the desired state.
// It uses early returns and a helper to reduce complexity.
func For(executedCommand, href string, options ...WaitOption) error {
	if href == "" {
		return nil
	}

	waitOpts := &WaitOptions{
		PollInterval: 3 * time.Second,
		Ctx:          context.Background(),
	}
	for _, opt := range options {
		opt(waitOpts)
	}

	c := client.Must().HttpClient
	for {
		// Check if the context is done.
		select {
		case <-waitOpts.Ctx.Done():
			return fmt.Errorf("timeout reached waiting for %s: %w", href, waitOpts.Ctx.Err())
		default:
		}

		resp, err := c.Get(href)
		if err != nil {
			// For delete, an error is considered a success (resource not found).
			if executedCommand == "delete" {
				return nil
			}
			return fmt.Errorf("failed to call %s: %w", href, err)
		}

		met, err := isConditionMet(executedCommand, resp)
		if err != nil {
			return err
		}
		if met {
			return nil
		}

		// Wait for the poll interval, checking context cancellation.
		select {
		case <-waitOpts.Ctx.Done():
			return fmt.Errorf("timeout reached waiting for %s: %w", href, waitOpts.Ctx.Err())
		case <-time.After(waitOpts.PollInterval):
		}
	}
}

// isConditionMet checks if the polling condition is met.
// For "delete": returns true if the HTTP status is 404.
// For "create" and "update": returns true if JSON's metadata.state is "AVAILABLE".
func isConditionMet(executedCommand string, resp *http.Response) (bool, error) {
	if executedCommand == "delete" {
		resp.Body.Close()
		return resp.StatusCode == http.StatusNotFound, nil
	}

	// Resource might not be available yet.
	if resp.StatusCode == 404 {
		return false, nil
	}

	// Error out if a 400 code.
	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("failed to call %s: %s", resp.Request.URL, resp.Status)
	}

	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode JSON: %w", err)
	}

	meta, ok := result["metadata"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("failed to parse metadata from response")
	}
	state, ok := meta["state"].(string)
	if !ok {
		return false, fmt.Errorf("failed to parse state from metadata")
	}
	return state == "AVAILABLE", nil
}
