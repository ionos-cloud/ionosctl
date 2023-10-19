package groups

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	vmasc "github.com/ionos-cloud/sdk-go-vmautoscaling"
	"github.com/spf13/viper"
)

func GroupCreateCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create VM Autoscaling Groups",
		Example: fmt.Sprintf("ionosctl vm-autoscaling group create %s",
			core.FlagsUsage(constants.FlagDatacenterId, constants.FlagName)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagDatacenterId)
			return err
		},
		CmdRun: func(c *core.CommandConfig) error {
			group, _, err := client.Must().VMAscClient.GroupsPost(context.Background()).GroupPost(vmasc.GroupPost{
				Properties: &vmasc.GroupProperties{
					Datacenter: &vmasc.GroupPropertiesDatacenter{
						Id: pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))),
					},
					Name: pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagName))),
				},
			}).Execute()
			if err != nil {
				return err
			}

			colsDesired := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))
			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, group,
				tabheaders.GetHeaders(allCols, defaultCols, colsDesired))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your group", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The ID of the datacenter")

	return cmd
}
