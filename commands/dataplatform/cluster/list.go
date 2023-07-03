package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
)

func ClusterListCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "dataplatform",
			Resource:  "cluster",
			Verb:      "list",
			Aliases:   []string{"l", "ls"},
			ShortDesc: "List Dataplatform Clusters",
			Example:   "ionosctl dataplatform cluster list",
			PreCmdRun: core.NoPreRun,
			CmdRun: func(c *core.CommandConfig) error {
				c.Printer.Verbose("Getting Clusters...")
				name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

				req := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(c.Context)
				if name != "" {
					req = req.Name(name)
				}
				clusters, _, err := req.Execute()
				if err != nil {
					return err
				}
				return c.Printer.Print(getClusterPrint(c, clusters.GetItems()))
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "",
		"Response filter to list only the clusters which include the specified name. case insensitive",
	)

	return cmd
}
