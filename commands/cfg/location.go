package cfg

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func LocationCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "location",
		Aliases:   []string{"location", "loc"},
		ShortDesc: "Print your config file's path",
		Example:   "ionosctl cfg location",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			cl, authErr := client.Get()
			var path string
			if authErr == nil && cl != nil && cl.ConfigPath != "" {
				path = cl.ConfigPath
			} else {
				path = config.GetConfigFilePath()
			}

			_, err := fmt.Fprintln(c.Command.Command.OutOrStdout(), path)
			return err
		},
		InitClient: false,
	})

	return cmd
}
