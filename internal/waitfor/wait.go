package waitfor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	core2 "github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

const (
	failed   = "FAILED"
	done     = "DONE"
	pollTime = 5 * time.Second
)

const (
	stateProgressCircleTpl   = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".. " "..." "...." }}`
	deleteProgressCircleTpl  = `{{ etime . }} {{ "Waiting for deletion" }}{{ cycle . "." ".. " "..." "...." }}`
	requestProgressCircleTpl = `{{ etime . }} {{ "Waiting for request" }}{{ cycle . "." ".. " "..." "...." }}`
)

type InterrogateRequestFunc func(c *core2.CommandConfig, requestId string) (status *string, message *string, err error)

// shouldWait checks whether waiting is requested, either via the per-command
// flag (e.g. --wait-for-request) or via the global --wait flag.
func shouldWait(c *core2.CommandConfig, perCmdFlag string) bool {
	if viper.GetBool(core2.GetFlagName(c.NS, perCmdFlag)) {
		return true
	}
	return viper.GetBool(constants.ArgWait)
}

// resolveTimeout returns the appropriate timeout: per-command --timeout if set,
// otherwise global --wait-timeout, falling back to DefaultTimeoutSeconds.
func resolveTimeout(c *core2.CommandConfig) int {
	if t := viper.GetInt(core2.GetFlagName(c.NS, constants.ArgTimeout)); t > 0 {
		return t
	}
	if t := viper.GetInt(constants.ArgWaitTimeout); t > 0 {
		return t
	}
	return constants.DefaultTimeoutSeconds
}

// WaitForRequest waits for Request to be executed
func WaitForRequest(c *core2.CommandConfig, interrogator InterrogateRequestFunc, requestId string) error {
	if !shouldWait(c, constants.ArgWaitForRequest) {
		return nil
	} else {
		timeout := resolveTimeout(c)
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		if isStructuredOutput() {
			return waitWithJSONLog(ctxTimeout, c, "Waiting for request...", func() <-chan error {
				_, errCh := WatchRequestProgress(ctxTimeout, c, interrogator, requestId)
				return errCh
			})
		}

		return waitWithProgressBar(c, requestProgressCircleTpl, func() <-chan error {
			_, errCh := WatchRequestProgress(ctxTimeout, c, interrogator, requestId)
			return errCh
		})
	}
}

type InterrogateStateFunc func(c *core2.CommandConfig, resourceId string) (*string, error)

func WaitForState(c *core2.CommandConfig, interrogator InterrogateStateFunc, resourceId string) error {
	if !shouldWait(c, constants.ArgWaitForState) {
		return nil
	} else {
		timeout := resolveTimeout(c)
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		if isStructuredOutput() {
			return waitWithJSONLog(ctxTimeout, c, "Waiting for state...", func() <-chan error {
				_, errCh := WatchStateProgress(ctxTimeout, c, interrogator, resourceId)
				return errCh
			})
		}

		return waitWithProgressBar(c, stateProgressCircleTpl, func() <-chan error {
			_, errCh := WatchStateProgress(ctxTimeout, c, interrogator, resourceId)
			return errCh
		})
	}
}

type InterrogateDeletionFunc func(c *core2.CommandConfig, resourceId string) (*int, error)

func WaitForDelete(c *core2.CommandConfig, interrogator InterrogateDeletionFunc, resourceId string) error {
	if !shouldWait(c, constants.ArgWaitForDelete) {
		return nil
	} else {
		timeout := resolveTimeout(c)
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		if isStructuredOutput() {
			return waitWithJSONLog(ctxTimeout, c, "Waiting for deletion...", func() <-chan error {
				_, errCh := WatchDeletionProgress(ctxTimeout, c, interrogator, resourceId)
				return errCh
			})
		}

		return waitWithProgressBar(c, deleteProgressCircleTpl, func() <-chan error {
			_, errCh := WatchDeletionProgress(ctxTimeout, c, interrogator, resourceId)
			return errCh
		})
	}
}

// isStructuredOutput returns true when the output format is JSON or api-json.
func isStructuredOutput() bool {
	switch viper.GetString(constants.ArgOutput) {
	case "json", "api-json":
		return true
	default:
		return false
	}
}

// logJSON writes a JSON-encoded string to stderr. This ensures the message is
// valid JSON so that tools like jq can skip it when processing a stream.
func logJSON(c *core2.CommandConfig, msg string) {
	out, _ := json.Marshal(msg)
	fmt.Fprintln(c.Command.Command.ErrOrStderr(), string(out))
}

// waitWithJSONLog emits JSON-formatted status messages to stderr (compatible
// with jq stream parsing) instead of an animated progress bar.
func waitWithJSONLog(_ context.Context, c *core2.CommandConfig, msg string, start func() <-chan error) error {
	logJSON(c, msg)
	if err := <-start(); err != nil {
		logJSON(c, failed)
		return err
	}
	logJSON(c, done)
	return nil
}

// waitWithProgressBar shows an animated progress bar on stderr for text output.
func waitWithProgressBar(c *core2.CommandConfig, tpl string, start func() <-chan error) error {
	progress := pb.New(1)
	progress.SetWriter(c.Command.Command.ErrOrStderr())
	progress.SetTemplateString(tpl)
	progress.Start()
	defer progress.Finish()

	if err := <-start(); err != nil {
		progress.SetTemplateString(tpl + " " + failed)
		return err
	}
	progress.SetTemplateString(tpl + " " + done)
	return nil
}
