package version

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
)

func Active() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "version",
		Verb:      "active",
		Aliases:   []string{"latest", "last"},
		ShortDesc: "Returns latest Dataplatform Cluster version, which you can use in cluster creation.",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateRawOutput(Latest(Versions())))
			return nil
		},
		InitClient: true,
	})

	return cmd
}
