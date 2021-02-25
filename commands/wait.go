package commands

import (
	"context"
	"fmt"
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

		// Default timeout: 60s
		timeout := viper.GetInt(builder.GetFlagName(c.ParentName, c.Name, config.ArgTimeout))
		ctxTimeout, cancel := context.WithTimeout(
			c.Context,
			time.Duration(timeout)*time.Second,
		)
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
