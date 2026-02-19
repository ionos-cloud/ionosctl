package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func GroupUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "group",
		Resource:  "group",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Group",
		LongDesc: `Use this command to update details about a specific Group.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Group Id`,
		Example:    "ionosctl compute group update --group-id GROUP_ID --reserve-ip",
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name for the Group")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateDc, "", false, "The group will be allowed to create Data Centers. E.g.: --create-dc=true, --create-dc=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateSnapshot, "", false, "The group will be allowed to create Snapshots. E.g.: --create-snapshot=true, --create-snapshot=false")
	cmd.AddBoolFlag(cloudapiv6.ArgReserveIp, "", false, "The group will be allowed to reserve IP addresses. E.g.: --reserve-ip=true, --reserve-ip=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessLog, "", false, "The group will be allowed to access the activity log. E.g.: --access-logs=true, --access-logs=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreatePcc, "", false, "The group will be allowed to create PCCs. E.g.: --create-pcc=true, --create-pcc=false")
	cmd.AddBoolFlag(cloudapiv6.ArgS3Privilege, "", false, "The group will be allowed to manage S3. E.g.: --s3privilege=true, --s3privilege=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateBackUpUnit, "", false, "The group will be able to manage Backup Units. E.g.: --create-backup=true, --create-backup=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateNic, "", false, "The group will be allowed to create NICs. E.g.: --create-nic=true, --create-nic=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateK8s, "", false, "The group will be allowed to create K8s Clusters. E.g.: --create-k8s=true, --create-k8s=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateFlowLog, "", false, "The group will be allowed to create Flow Logs. E.g.: --create-flowlog=true, --create-flowlog=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessMonitoring, "", false, "Privilege for a group to access and manage monitoring related functionality using Monotoring-as-a-Service. E.g.: --access-monitoring=true, --access-monitoring=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessCerts, "", false, "Privilege for a group to access and manage certificates. E.g.: --access-certs=true, --access-certs=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessDNS, "", false, "Privilege for a group to access and manage dns records")
	cmd.AddBoolFlag(cloudapiv6.ArgManageDbaas, "", false, "Privilege for a group to manage DBaaS related functionality")
	cmd.AddBoolFlag(cloudapiv6.ArgManageRegistry, "", false, "Privilege for group accessing container registry related functionality")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for Group update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Group update [seconds]")

	return cmd
}
