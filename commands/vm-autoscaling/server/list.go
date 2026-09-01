package server

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
		ShortDesc: "List the replica VMs currently running in one or all groups",
		LongDesc:  "List the replica VMs the autoscaler is currently running. Pass --group-id to see the replicas of one group, or --all to gather replicas across every group in your account (this fetches servers group-by-group, so it is slower with many groups). Each row is enriched with the underlying Compute Engine (CloudAPI) server's live details - name, availability zone, cores, RAM, CPU family and state - so the list reflects the actual running VMs, not just their IDs.",
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

			return c.Printer(allCols).Print(enriched)
		},
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List replicas from every VM Auto Scaling group in your account (mutually exclusive with --group-id)")
	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "The ID of the VM Auto Scaling group whose replica VMs to list")
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

	return c.Printer(allCols).Print(enriched)
}
