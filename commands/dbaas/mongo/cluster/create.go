package cluster

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"
	"github.com/spf13/cobra"
)

func ClusterCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Aliases:   []string{"c"},
		ShortDesc: "Create Mongo Clusters",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Creating Cluster...")
			clusters, _, err := c.DbaasMongoServices.Clusters().List("")
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(nil, c, clusters.GetItems()))
		},
		InitClient: true,
	})

	// TODO: Move ArgName to DBAAS level constants
	cmd.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	cmd.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(config.ArgCols, "", allCols[0:6], printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
