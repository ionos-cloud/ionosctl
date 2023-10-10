package cfg

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/shellpath"
)

func LocationCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "location",
		Aliases:   []string{"location", "loc"},
		ShortDesc: "Print your config file's path",
		Example:   "ionosctl cfg loc",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			_, err := fmt.Fprintf(c.Command.Command.OutOrStdout(), shellpath.EscapePathForShell(config.GetConfigFilePath()))
			return err
		},
		InitClient: false,
	})

	return cmd
}
