package group

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "group",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a VM Autoscaling Group",
		Example: fmt.Sprintf("ionosctl vm-autoscaling group get %s",
			core.FlagsUsage(constants.FlagGroupId)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagGroupId, constants.FlagServerId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().VMAscClient.GroupsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId))).
				Depth(float32(viper.GetFloat64(core.GetFlagName(c.NS, constants.ArgDepth)))).
				Execute()
			if err != nil {
				return err
			}

			colsDesired := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.VmAutoscalingGroup, ls,
				tabheaders.GetHeaders(allCols, defaultCols, colsDesired))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddInt32Flag(constants.ArgDepth, constants.ArgDepthShort, 1, "Controls the detail depth of the response objects")
	cmd.AddStringFlag(constants.FlagGroupId, "", "", "ID of the autoscaling group")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return GroupsProperty(func(r vmasc.Group) string {
			completion := *r.Id
			if r.Properties.Name == nil {
				return completion
			}
			completion += "\t" + *r.Properties.Name
			return completion
		}, func(r vmasc.ApiGroupsGetRequest) (vmasc.ApiGroupsGetRequest, error) {
			return r.Depth(1), nil
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
