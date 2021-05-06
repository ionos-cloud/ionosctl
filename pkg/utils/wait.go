package utils

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
)

const (
	failed    = "FAILED"
	active    = "ACTIVE"
	available = "AVAILABLE"
	pollTime  = 10
)

var waitingForRequestMsg = "Waiting for request: %s"
var waitingForStateMsg = "Waiting for state: %s"
var contextTimeoutErr = errors.New("context hit timeout")

// WaitForRequest waits for request to be executed
func WaitForRequest(c *builder.CommandConfig, path string) error {
	if !viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)) {
		return nil
	} else {
		timeout := viper.GetInt(builder.GetFlagName(c.ParentName, c.Name, config.ArgTimeout))
		if timeout == 0 {
			timeout = config.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()
		c.Context = ctxTimeout

		reqId, err := printer.GetRequestId(path)
		if err != nil {
			return err
		}
		if err = c.Printer.Print(fmt.Sprintf(waitingForRequestMsg, *reqId)); err != nil {
			return err
		}
		if _, err = c.Requests().Wait(path); err != nil {
			return err
		}
	}
	return nil
}

type InterrogateStateFunc func(c *builder.CommandConfig, resourceId string) (string, error)

// WaitForState waits for the State of a Resource to be Active or Available
func WaitForState(c *builder.CommandConfig, interrog InterrogateStateFunc, resourceId string) error {
	if !viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForState)) {
		return nil
	} else {
		timeout := viper.GetInt(builder.GetFlagName(c.ParentName, c.Name, config.ArgTimeout))
		if timeout == 0 {
			timeout = config.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)

		var wg sync.WaitGroup
		wg.Add(1)
		go func(cmdCfg *builder.CommandConfig, interrogator InterrogateStateFunc, resId string) {
			for {
				select {
				case <-ctxTimeout.Done():
					wg.Done()
					clierror.CheckError(contextTimeoutErr, cmdCfg.Printer.GetStderr())
					return
				default:
					if state, err := interrogator(cmdCfg, resId); err == nil {
						cmdCfg.Printer.Print(fmt.Sprintf(waitingForStateMsg, state))
						if IsActive(state) {
							cmdCfg.Printer.Print(state)
							wg.Done()
							return
						}
						if HasFailed(state) {
							cmdCfg.Printer.Print(state)
							wg.Done()
							return
						}
					}

					time.Sleep(pollTime * time.Second)
				}
			}
		}(c, interrog, resourceId)
		wg.Wait()
		cancel()
	}
	return nil
}

func IsActive(state string) bool {
	if state == active || state == available {
		return true
	} else {
		return false
	}
}

func HasFailed(state string) bool {
	if state == failed {
		return true
	} else {
		return false
	}
}
