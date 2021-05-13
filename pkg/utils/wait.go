package utils

import (
	"context"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
)

const (
	failed   = "FAILED"
	done     = "DONE"
	pollTime = 10
)

const (
	stateProgressCircleTpl   = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".. " "..." "...." }}`
	requestProgressCircleTpl = `{{ etime . }} {{ "Waiting for request" }}{{ cycle . "." ".. " "..." "...." }}`
)

var waitingForRequestMsg = "Waiting for request..."
var waitingForStateMsg = "Waiting for state..."

// WaitForRequest waits for Request to be executed
func WaitForRequest(c *core.CommandConfig, requestPath string) error {
	if !viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest)) {
		// Double Check: return if flag not set
		return nil
	} else {
		// Set context timeout
		timeout := viper.GetInt(core.GetFlagName(c.NS, config.ArgTimeout))
		if timeout == 0 {
			timeout = config.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		// Get Request Id
		requestId, err := printer.GetRequestId(requestPath)
		if err != nil {
			return err
		}

		// Check the output format
		if viper.GetString(config.ArgOutput) == printer.TypeText.String() {
			progress := pb.New(1)
			progress.SetWriter(c.Printer.GetStdout())
			progress.SetTemplateString(requestProgressCircleTpl)
			progress.Start()
			defer progress.Finish()

			_, errCh := WatchRequestProgress(ctxTimeout, c, *requestId)
			if err := <-errCh; err != nil {
				progress.SetTemplateString(requestProgressCircleTpl + " " + failed)
				return err
			}
			progress.SetTemplateString(requestProgressCircleTpl + " " + done)
		} else {
			c.Printer.Print(waitingForRequestMsg)
			_, errCh := WatchRequestProgress(ctxTimeout, c, *requestId)
			if err := <-errCh; err != nil {
				c.Printer.Print(failed)
				return err
			}
			c.Printer.Print(done)
		}
		return nil
	}
}

type InterrogateStateFunc func(c *core.CommandConfig, resourceId string) (*string, error)

func WaitForState(c *core.CommandConfig, interrogator InterrogateStateFunc, resourceId string) error {
	if !viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		// Double Check: return if flag not set
		return nil
	} else {
		// Set context timeout
		timeout := viper.GetInt(core.GetFlagName(c.NS, config.ArgTimeout))
		if timeout == 0 {
			timeout = config.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		// Check the output format
		if viper.GetString(config.ArgOutput) == printer.TypeText.String() {
			progress := pb.New(1)
			progress.SetWriter(c.Printer.GetStdout())
			progress.SetTemplateString(stateProgressCircleTpl)
			progress.Start()
			defer progress.Finish()

			_, errCh := WatchStateProgress(ctxTimeout, c, interrogator, resourceId)
			if err := <-errCh; err != nil {
				progress.SetTemplateString(stateProgressCircleTpl + " " + failed)
				return err
			}
			progress.SetTemplateString(stateProgressCircleTpl + " " + done)
		} else {
			c.Printer.Print(waitingForStateMsg)
			_, errCh := WatchStateProgress(ctxTimeout, c, interrogator, resourceId)
			if err := <-errCh; err != nil {
				c.Printer.Print(failed)
				return err
			}
			c.Printer.Print(done)
		}
		return nil
	}
}
