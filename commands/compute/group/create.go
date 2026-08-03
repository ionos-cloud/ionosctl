package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func GroupCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "group",
		Resource:  "group",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Group with a set of privileges",
		LongDesc: `Create a new IAM Group. A Group bundles a set of contract-wide PRIVILEGES; every User you later add to the Group (via ` + "`ionosctl compute group user add`" + `) inherits all of them.

Each privilege flag is an independent capability and defaults to false, so a freshly created Group grants nothing until you enable the flags you need. Privileges are contract-wide "may do X anywhere" capabilities - they are separate from resource shares (` + "`ionosctl compute share`" + `), which grant access to one specific resource.

Only the name is truly required; pass any combination of privilege flags to switch capabilities on at creation time.`,
		Example: `# Minimal group (no privileges yet)
ionosctl compute group create --name "Developers"

# A team that can spin up infrastructure: datacenters, K8s clusters, snapshots and public connectivity
ionosctl compute group create --name "Platform" --create-dc --create-k8s --create-snapshot --create-nic

# A read-only auditing group that can only inspect the activity log
ionosctl compute group create --name "Auditors" --access-logs`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunGroupCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Group", "The name of the Group. Not required to be unique, but used to identify the Group in listings and shares")
	addGroupPrivilegeFlags(cmd)

	return cmd
}

// addGroupPrivilegeFlags registers the identical set of contract-wide privilege
// toggles shared by `group create` and `group update`. Each flag grants (or, on
// update, revokes) one capability to every member of the Group. Descriptions map
// each CLI flag to the underlying IAM privilege it controls, since a few flag
// names differ from the privilege they set (e.g. --create-nic controls the
// "create internet access" privilege, --access-logs the "access activity log"
// privilege).
func addGroupPrivilegeFlags(cmd *core.Command) {
	cmd.AddBoolFlag(cloudapiv6.ArgCreateDc, "", false, "Grant the 'create data center' privilege: members may create new Virtual Data Centers. E.g.: --create-dc=true, --create-dc=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateSnapshot, "", false, "Grant the 'create snapshot' privilege: members may take snapshots of volumes. E.g.: --create-snapshot=true, --create-snapshot=false")
	cmd.AddBoolFlag(cloudapiv6.ArgReserveIp, "", false, "Grant the 'reserve IP' privilege: members may reserve public IP blocks. E.g.: --reserve-ip=true, --reserve-ip=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessLog, "", false, "Grant the 'access activity log' privilege: members may read the contract's audit/activity log. E.g.: --access-logs=true, --access-logs=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreatePcc, "", false, "Grant the 'create PCC' privilege: members may create Private Cross-Connects to bridge private LANs across datacenters. E.g.: --create-pcc=true, --create-pcc=false")
	cmd.AddBoolFlag(cloudapiv6.ArgS3Privilege, "", false, "Grant the S3 privilege: members may use IONOS Object Storage (S3-compatible) and manage their own S3 keys. E.g.: --s3privilege=true, --s3privilege=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateBackUpUnit, "", false, "Grant the 'create backup unit' privilege: members may create and manage Backup Units. E.g.: --create-backup=true, --create-backup=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateNic, "", false, "Grant the 'create internet access' privilege: members may attach public/internet-facing connectivity (despite the flag name, this is NOT a per-NIC toggle). E.g.: --create-nic=true, --create-nic=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateK8s, "", false, "Grant the 'create Kubernetes cluster' privilege: members may create Managed Kubernetes clusters. E.g.: --create-k8s=true, --create-k8s=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCreateFlowLog, "", false, "Grant the 'create Flow Log' privilege: members may create Flow Logs to capture network traffic. E.g.: --create-flowlog=true, --create-flowlog=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessMonitoring, "", false, "Grant the 'access and manage monitoring' privilege: members may access metrics and manage alarms via Monitoring-as-a-Service (MaaS). E.g.: --access-monitoring=true, --access-monitoring=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessCerts, "", false, "Grant the 'access and manage certificates' privilege: members may manage certificates in the Certificate Manager. E.g.: --access-certs=true, --access-certs=false")
	cmd.AddBoolFlag(cloudapiv6.ArgAccessDNS, "", false, "Grant the 'access and manage DNS' privilege: members may manage DNS zones and records")
	cmd.AddBoolFlag(cloudapiv6.ArgManageDbaas, "", false, "Grant the 'manage DBaaS' privilege: members may manage Database-as-a-Service clusters (PostgreSQL, MongoDB, MariaDB, etc.)")
	cmd.AddBoolFlag(cloudapiv6.ArgManageRegistry, "", false, "Grant the 'manage Registry' privilege: members may manage Container Registry repositories")
}
