package waitfor

import (
	"context"
	"fmt"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	core2 "github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/spf13/viper"
)

const (
	failed   = "FAILED"
	done     = "DONE"
	pollTime = time.Second
)

const (
	stateProgressCircleTpl   = `{{ etime . }} {{ "Waiting for state" }}{{ cycle . "." ".. " "..." "...." }}`
	deleteProgressCircleTpl  = `{{ etime . }} {{ "Waiting for deletion" }}{{ cycle . "." ".. " "..." "...." }}`
	requestProgressCircleTpl = `{{ etime . }} {{ "Waiting for request" }}{{ cycle . "." ".. " "..." "...." }}`
)

var waitingForRequestMsg = "Waiting for request..."
var waitingForStateMsg = "Waiting for state..."

type InterrogateRequestFunc func(c *core2.CommandConfig, requestId string) (status *string, message *string, err error)

// WaitForRequest waits for Request to be executed
func WaitForRequest(c *core2.CommandConfig, interrogator InterrogateRequestFunc, requestId string) error {
	if !viper.GetBool(core2.GetFlagName(c.NS, constants.ArgWaitForRequest)) {
		// Double Check: return if flag not set
		return nil
	} else {
		// Set context timeout
		timeout := viper.GetInt(core2.GetFlagName(c.NS, constants.ArgTimeout))
		if timeout == 0 {
			timeout = constants.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		// Check the output format
		if viper.GetString(constants.ArgOutput) == jsontabwriter.TextFormat {
			progress := pb.New(1)
			progress.SetWriter(c.Command.Command.OutOrStdout())
			progress.SetTemplateString(requestProgressCircleTpl)
			progress.Start()
			defer progress.Finish()

			_, errCh := WatchRequestProgress(ctxTimeout, c, interrogator, requestId)
			if err := <-errCh; err != nil {
				progress.SetTemplateString(requestProgressCircleTpl + " " + failed)
				return err
			}
			progress.SetTemplateString(requestProgressCircleTpl + " " + done)
		} else {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(waitingForRequestMsg))
			_, errCh := WatchRequestProgress(ctxTimeout, c, interrogator, requestId)
			if err := <-errCh; err != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(failed))
				return err
			}
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(done))
		}
		return nil
	}
}

type InterrogateStateFunc func(c *core2.CommandConfig, resourceId string) (*string, error)

func WaitForState(c *core2.CommandConfig, interrogator InterrogateStateFunc, resourceId string) error {
	if !viper.GetBool(core2.GetFlagName(c.NS, constants.ArgWaitForState)) {
		// Double Check: return if flag not set
		return nil
	} else {
		// Set context timeout
		timeout := viper.GetInt(core2.GetFlagName(c.NS, constants.ArgTimeout))
		if timeout == 0 {
			timeout = constants.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		// Check the output format
		if viper.GetString(constants.ArgOutput) == jsontabwriter.TextFormat {
			progress := pb.New(1)
			progress.SetWriter(c.Command.Command.OutOrStdout())
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
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(waitingForStateMsg))
			_, errCh := WatchStateProgress(ctxTimeout, c, interrogator, resourceId)
			if err := <-errCh; err != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(failed))
				return err
			}
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(done))
		}
		return nil
	}
}

type InterrogateDeletionFunc func(c *core2.CommandConfig, resourceId string) (*int, error)

func WaitForDelete(c *core2.CommandConfig, interrogator InterrogateDeletionFunc, resourceId string) error {
	if !viper.GetBool(core2.GetFlagName(c.NS, constants.ArgWaitForDelete)) {
		// Double Check: return if flag not set
		return nil
	} else {
		// Set context timeout
		timeout := viper.GetInt(core2.GetFlagName(c.NS, constants.ArgTimeout))
		if timeout == 0 {
			timeout = constants.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		// Check the output format
		if viper.GetString(constants.ArgOutput) == jsontabwriter.TextFormat {
			progress := pb.New(1)
			progress.SetWriter(c.Command.Command.OutOrStdout())
			progress.SetTemplateString(deleteProgressCircleTpl)
			progress.Start()
			defer progress.Finish()

			// WaitForDelete monitors the http Response Status Code.
			_, errCh := WatchDeletionProgress(ctxTimeout, c, interrogator, resourceId)
			if err := <-errCh; err != nil {
				progress.SetTemplateString(deleteProgressCircleTpl + " " + failed)
				return err
			}
			progress.SetTemplateString(deleteProgressCircleTpl + " " + done)
		} else {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(waitingForStateMsg))
			_, errCh := WatchDeletionProgress(ctxTimeout, c, interrogator, resourceId)
			if err := <-errCh; err != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(failed))
				return err
			}
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(done))
		}
		return nil
	}
}
