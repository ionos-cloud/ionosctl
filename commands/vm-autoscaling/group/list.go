package group

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

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List VM Autoscaling Groups. Use a greater '--depth' to see current replica count",
		Example:   "ionosctl vm-autoscaling group list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().VMAscClient.GroupsGet(context.Background()).
				Depth(float32(viper.GetFloat64(core.GetFlagName(c.NS, constants.ArgDepth)))).Execute()
			if err != nil {
				return err
			}

			table, err := resource2table.ConvertVmAutoscalingGroupsToTable(ls)
			if err != nil {
				return err
			}

			colsDesired := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))
			out, err := jsontabwriter.GenerateOutputPreconverted(ls, table,
				tabheaders.GetHeaders(allCols, defaultCols, colsDesired))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddInt32Flag(constants.ArgDepth, constants.ArgDepthShort, 1, "Controls the detail depth of the response objects")

	return cmd
}
