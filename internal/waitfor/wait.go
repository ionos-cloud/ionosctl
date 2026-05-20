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
	stateProgressCircleTpl = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".. " "..." "...." }}`
)

type InterrogateStateFunc func(c *core2.CommandConfig, resourceId string) (*string, error)

// Deprecated: WaitForState is a legacy wait function. New code should rely on
// globalwait.WaitAndRerender which handles waiting automatically via --wait.
// Only remaining caller: commands/compute/server/run_server.go (--promote-volume).
func WaitForState(c *core2.CommandConfig, interrogator InterrogateStateFunc, resourceId string) error {
	if !viper.GetBool(constants.ArgWait) {
		return nil
	} else {
		timeout := viper.GetInt(constants.ArgTimeout)
		if timeout <= 0 {
			timeout = constants.DefaultTimeoutSeconds
		}
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
