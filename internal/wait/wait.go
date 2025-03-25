package wait

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
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
func For(commandName, href string, options ...WaitOption) {
	fmt.Fprint(os.Stderr, jsontabwriter.GenerateVerboseOutput("Waiting for "+href+" to reach desired state..."))

	waitOpts := &WaitOptions{
		PollInterval: 3 * time.Second,
		Ctx:          context.Background(),
	}
	for _, opt := range options {
		opt(waitOpts)
	}

	c := client.Must().HttpClient
	for {
		select {
		case <-waitOpts.Ctx.Done():
			fmt.Fprint(os.Stderr, jsontabwriter.GenerateVerboseOutput("Stopped waiting: context timeout"))
			return
		default:
		}

		resp, err := c.Get(href)
		if err != nil {
			// if commandName == "delete" && resp != nil && resp.StatusCode == http.StatusNotFound {
			// 	fmt.Fprintln(os.Stderr, jsontabwriter.GenerateVerboseOutput("Resource successfully deleted"))
			// 	return
			// }
			fmt.Fprint(os.Stderr, jsontabwriter.GenerateVerboseOutput("Failed to call "+href+": "+err.Error()))
			return // currently simply stop waiting if an error occurs
		}

		met, err := isConditionMet(commandName, resp)
		if err != nil {
			fmt.Fprint(os.Stderr, jsontabwriter.GenerateVerboseOutput("Failed waiting for condition to be met: "+err.Error()))
			return // currently simply stop waiting if an error occurs
		}
		if met {
			fmt.Fprint(os.Stderr, jsontabwriter.GenerateVerboseOutput("Successfully waited for condition to be met"))
			return
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
