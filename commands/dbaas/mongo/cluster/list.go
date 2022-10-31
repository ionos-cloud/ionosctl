package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"
	"github.com/spf13/cobra"
)

func ClusterListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Mongo Clusters",
		LongDesc:  "Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:   "ionosctl dbaas mongo cluster list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Clusters...")
			clusters, _, err := c.DbaasMongoServices.Clusters().List("")
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(c, clusters.GetItems()))
		},
		InitClient: true,
	})

	// TODO: Move ArgName to DBAAS level constants
	cmd.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", allCols[0:6], printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
