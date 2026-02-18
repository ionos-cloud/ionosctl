package targetgroup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultTargetGroupTargetCols = []string{"TargetIp", "TargetPort", "Weight", "HealthCheckEnabled", "MaintenanceEnabled"}
)

func TargetGroupTargetCmd() *core.Command {
	targetGroupTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "target",
			Aliases:          []string{"t"},
			Short:            "Target Group Target Operations",
			Long:             "The sub-commands of `ionosctl targetgroup target` allow you to see information, to add, remove Targets from Target Groups.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(context.TODO(), targetGroupTargetCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "target",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Target Groups Targets",
		LongDesc:   "Use this command to get a list of Target Groups Targets.",
		Example:    `ionosctl targetgroup target list --targetgroup-id TARGET_GROUP_ID`,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupTargetList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.ArgCols, "", defaultTargetGroupTargetCols, tabheaders.ColsMessage(defaultTargetGroupTargetCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTargetGroupTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	add := core.NewCommand(context.TODO(), targetGroupTargetCmd, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "target",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Target to a Target Group",
		LongDesc: `Use this command to add a Target to a Target Group. You will need to provide the IP, the port and the weight. Weight parameter is used to adjust the target VM's weight relative to other target VMs. All target VMs will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. The default weight is 1, and the maximal value is 256. A value of 0 means the target VM will not participate in load-balancing but will still accept persistent connections. If this parameter is used to distribute the load according to target VM's capacity, it is recommended to start with values which can both grow and shrink, for instance between 10 and 100 to leave enough room above and below for later adjustments.

Health Check can also be set. The ` + "`" + `--check` + "`" + ` option specifies whether the target VM's health is checked. If turned off, a target VM is always considered available. If turned on, the target VM is available when accepting periodic TCP connections, to ensure that it is really able to serve requests. The address and port to send the tests to are those of the target VM. The health check only consists of a connection attempt.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port`,
		Example:    `ionosctl targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT`,
		PreCmdRun:  PreRunTargetGroupIdTargetIpPort,
		CmdRun:     RunTargetGroupTargetAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The IP of the balanced target VM.", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, 8080, "The port of the balanced target service; valid range is 1 to 65535.", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgWeight, cloudapiv6.ArgWeightShort, 1, "Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments.")
	add.AddBoolFlag(cloudapiv6.ArgHealthCheckEnabled, "", true, "Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True.")
	add.AddBoolFlag(cloudapiv6.ArgMaintenanceEnabled, cloudapiv6.ArgMaintenanceShort, false, "Maintenance mode prevents the target from receiving balanced traffic.")
	add.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group Target addition to be executed")
	add.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group Target addition [seconds]")
	add.AddStringSliceFlag(constants.ArgCols, "", defaultTargetGroupTargetCols, tabheaders.ColsMessage(defaultTargetGroupTargetCols))
	_ = add.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTargetGroupTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove Command
	*/
	remove := core.NewCommand(context.TODO(), targetGroupTargetCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "target",
		Verb:       "remove",
		Aliases:    []string{"r"},
		ShortDesc:  "Remove a Target from a Target Group",
		LongDesc:   "Use this command to delete the specified Target from Target Group.\n\nRequired values to run command:\n\n* Target Group Id\n* Target Ip\n* Target Port",
		Example:    `ionosctl targetgroup target remove --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT`,
		PreCmdRun:  PreRunTargetGroupTargetRemove,
		CmdRun:     RunTargetGroupTargetRemove,
		InitClient: true,
	})
	remove.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = remove.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	remove.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP of a balanced target VM", core.RequiredFlagOption())
	remove.AddIntFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, 8080, "Port of the balanced target service. (range: 1 to 65535)", core.RequiredFlagOption())
	remove.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Target Group Targets")
	remove.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group Target deletion to be executed")
	remove.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group Target deletion [seconds]")
	remove.AddStringSliceFlag(constants.ArgCols, "", defaultTargetGroupTargetCols, tabheaders.ColsMessage(defaultTargetGroupTargetCols))
	_ = remove.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTargetGroupTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(targetGroupTargetCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
