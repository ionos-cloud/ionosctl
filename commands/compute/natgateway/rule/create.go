package rule

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NatgatewayRuleCreateCmd() *core.Command {
	create := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NAT Gateway Rule",
		LongDesc: `Use this command to create a NAT Gateway Rule in a specified NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Public IP
* Source Subnet`,
		Example:    "ionosctl natgateway rule create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --ip IP_1 --source-subnet SOURCE_SUBNET --target-subnet TARGET_SUBNET",
		PreCmdRun:  cloudapiv6cmds.PreRunNatGatewayRuleCreate,
		CmdRun:     cloudapiv6cmds.RunNatGatewayRuleCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Rule", "Name of the NAT Gateway Rule")
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, string(ionoscloud.ALL), "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIpFlag(cloudapiv6.ArgIp, "", nil, "Public IP address of the NAT Gateway Rule", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	create.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	create.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule creation [seconds]")

	return create
}
