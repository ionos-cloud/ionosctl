package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func ClusterListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Dataplatform Clusters",
		Example:   "ionosctl dataplatform cluster list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Clusters..."))

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			req := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(c.Context)
			if name != "" {
				req = req.Name(name)
			}

			clusters, _, err := req.Execute()
			if err != nil {
				return err
			}

			clustersConverted, err := resource2table.ConvertDataplatformClustersToTable(clusters)
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutputPreconverted(clusters, clustersConverted,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the clusters which include the specified name. case insensitive")

	return cmd
}
