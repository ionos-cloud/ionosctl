package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func GroupCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "group",
		Verb:       "create",
		Aliases:    []string{"c"},
		ShortDesc:  "Create a Group",
		LongDesc:   `Use this command to create a new Group and set Group privileges. You can specify the name for the new Group. By default, all privileges will be set to false. You need to use flags privileges to be set to true.`,
		Example:    "ionosctl group create --name NAME --wait-for-request",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunGroupCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Group", "Name for the Group")
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
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for Group creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Group creation [seconds]")

	return cmd
}
