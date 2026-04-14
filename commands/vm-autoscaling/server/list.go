package server

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "server",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Servers that are managed by VM-Autoscaling Groups",
		Example: fmt.Sprintf(`ionosctl vm-autoscaling server list %s
ionosctl vm-autoscaling server list %s`,
			core.FlagUsage(constants.FlagGroupId), core.FlagUsage(constants.ArgAll)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGroupId},
				[]string{constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
				return listAll(c)
			}

			ls, _, err := client.Must().VMAscClient.GroupsServersGet(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId))).
				Execute()
			if err != nil {
				return err
			}

			enriched, err := enrichAutoscalingServers(ls)
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, enriched, cols))
		},
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "If set, list all servers of all groups")
	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "ID of the autoscaling group to list servers from")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// get ID of all groups
		return group.GroupsProperty(func(r vmasc.Group) string {
			completion := *r.Id
			if r.Properties == nil || r.Properties.Name == nil {
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

func listAll(c *core.CommandConfig) error {
	ls, err := Servers()
	if err != nil {
		return fmt.Errorf("failed listing servers of all groups: %w", err)
	}

	enriched, err := enrichAutoscalingServers(ls)
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, enriched, cols))
}

