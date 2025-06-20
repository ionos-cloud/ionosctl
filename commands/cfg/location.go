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
			path := client.Must().ConfigPath

			cl, authErr := client.Get()
			if cl == nil || authErr != nil || cl.Config == nil || path == "" {
				path = config.GetConfigFilePath() // fallback to default ionosctl config path
			}

			_, err := fmt.Fprintf(c.Command.Command.OutOrStdout(), path)
			return err
		},
		InitClient: false,
	})

	return cmd
}
