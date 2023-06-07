package cfg

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func CfgLocationCmd() *core.Command {
	loginCmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "config",
		Resource:  "config",
		Verb:      "config",
		Aliases:   []string{"cfg"},
		ShortDesc: "Print your config file's path",
		Example:   "ionosctl cfg",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			_, err := fmt.Fprintf(c.Command.Command.OutOrStdout(), config.GetConfigFile())
			return err
		},
		InitClient: false,
	})

	return loginCmd
}
