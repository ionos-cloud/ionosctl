package flowlog

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FlowLogCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "flowlog",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a FlowLog on a NIC",
		LongDesc: `Use this command to create a Flow Log on a specific NIC. The Flow Log streams traffic metadata for that NIC into an existing IONOS Object Storage bucket.

The bucket named by --s3bucket must already exist and be in a location that supports Flow Log delivery; this command does not create it. Choose which traffic to capture with --action (ACCEPTED / REJECTED / ALL) and --direction (INGRESS / EGRESS / BIDIRECTIONAL).

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

NOTE: Disable/delete the Flow Log before deleting the Object Storage bucket it writes to.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Target S3 Bucket Name`,
		Example: `# Capture all traffic in both directions
ionosctl compute flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name "audit-all" --s3bucket my-flowlog-bucket

# Capture only rejected inbound traffic (useful for security review) and wait until ready
ionosctl compute flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name "blocked-ingress" --action REJECTED --direction INGRESS --s3bucket my-flowlog-bucket --wait`,
		PreCmdRun:  PreRunFlowLogCreate,
		CmdRun:     RunFlowLogCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", "The unique ID of the Virtual Data Center that holds the server and NIC", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", "The unique ID of the server that owns the NIC", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", "The unique ID of the NIC whose traffic will be captured", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed FlowLog", "A display name for the Flow Log")
	cmd.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "ALL", "Which traffic to capture by firewall outcome: ACCEPTED (allowed traffic), REJECTED (blocked traffic), or ALL (both)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "BIDIRECTIONAL", "Which traffic direction to capture: INGRESS (inbound to the NIC), EGRESS (outbound from the NIC), or BIDIRECTIONAL (both)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "Name of an EXISTING IONOS Object Storage (S3) bucket that will receive the Flow Log files. The bucket must already exist", core.RequiredFlagOption())

	return cmd
}
