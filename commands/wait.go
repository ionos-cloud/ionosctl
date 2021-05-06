package commands

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
)

var (
	waitingForActionMsg = "Waiting for request: %s"
)

func waitForAction(c *builder.CommandConfig, path string) error {
	if !viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)) {
		return nil
	} else {
		reqId, err := printer.GetRequestId(path)
		if err != nil {
			return err
		}

		timeout := viper.GetInt(builder.GetFlagName(c.ParentName, c.Name, config.ArgTimeout))
		if timeout == 0 {
			timeout = config.DefaultTimeoutSeconds
		}
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)
		defer cancel()

		c.Context = ctxTimeout
		err = c.Printer.Print(fmt.Sprintf(waitingForActionMsg, *reqId))
		if err != nil {
			return err
		}
		if _, err = c.Requests().Wait(path); err != nil {
			return err
		}
	}
	return nil
}

type InterrogateStateFunc func(c *builder.CommandConfig, resourceId string) (string, error)

func IsActive(state string) bool {
	fmt.Println(state)
	if state == "ACTIVE" || state == "AVAILABLE" {
		return true
	} else {
		return false
	}
}

func waitForState(c *builder.CommandConfig, interrog InterrogateStateFunc, resourceId string) error {
	if !viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)) {
		return nil
	} else {
		timeout := viper.GetInt(builder.GetFlagName(c.ParentName, c.Name, config.ArgTimeout))
		if timeout == 0 {
			timeout = config.DefaultTimeoutSeconds
		}
		fmt.Println("timeout:")
		fmt.Println(timeout)
		ctxTimeout, cancel := context.WithTimeout(c.Context, time.Duration(timeout)*time.Second)

		var wg sync.WaitGroup
		wg.Add(1)
		go func(cmdCfg *builder.CommandConfig, interrogator InterrogateStateFunc, resId string) {
			for {
				select {
				case <-ctxTimeout.Done():
					cmdCfg.Printer.Print("DONE")
					wg.Done()
					return
				default:
					cmdCfg.Printer.Print("Waiting for state..")
					// interogate state if IsActive break
					if state, err := interrogator(cmdCfg, resId); err == nil && IsActive(state) {
						cmdCfg.Printer.Print("Done state")
						wg.Done()
						return
					}
					time.Sleep(5 * time.Second)
				}
			}
		}(c, interrog, resourceId)
		wg.Wait()
		cancel()
	}
	return nil
}
