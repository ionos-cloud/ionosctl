package cluster

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"os"
)

const (
	flagClusterId      = "cluster-id"
	flagClusterIdShort = "i"
)

func ClusterGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a Mongo Cluster by ID",
		Example:   "ionosctl dbaas mongo cluster get --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(flagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId, err := c.Command.Command.Flags().GetString(flagClusterId)
			if err != nil {
				return err
			}
			c.Printer.Verbose("Getting Cluster by id: %s", clusterId)
			cluster, r, err := c.DbaasMongoServices.Clusters().Get(clusterId)
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(r, c, &[]ionoscloud.ClusterResponse{cluster}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(flagClusterId, flagClusterIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(config.ArgCols, "", allCols[0:6], printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
