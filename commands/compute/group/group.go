package group

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allGroupCols = []table.Column{
	{Name: "GroupId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "CreateDataCenter", JSONPath: "properties.createDataCenter", Default: true},
	{Name: "CreateSnapshot", JSONPath: "properties.createSnapshot", Default: true},
	{Name: "CreatePcc", JSONPath: "properties.createPcc", Default: true},
	{Name: "CreateBackupUnit", JSONPath: "properties.createBackupUnit", Default: true},
	{Name: "CreateInternetAccess", JSONPath: "properties.createInternetAccess", Default: true},
	{Name: "CreateK8s", JSONPath: "properties.createK8sCluster", Default: true},
	{Name: "ReserveIp", JSONPath: "properties.reserveIp", Default: true},
	{Name: "AccessActivityLog", JSONPath: "properties.accessActivityLog"},
	{Name: "S3Privilege", JSONPath: "properties.s3Privilege"},
	{Name: "CreateFlowLog", JSONPath: "properties.createFlowLog"},
	{Name: "AccessAndManageMonitoring", JSONPath: "properties.accessAndManageMonitoring"},
	{Name: "AccessAndManageCertificates", JSONPath: "properties.accessAndManageCertificates"},
	{Name: "AccessAndManageDns", JSONPath: "properties.accessAndManageDns"},
	{Name: "ManageDBaaS", JSONPath: "properties.manageDBaaS"},
	{Name: "ManageRegistry", JSONPath: "properties.manageRegistry"},
}

func GroupCmd() *core.Command {
	groupCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "group",
			Aliases: []string{"g"},
			Short:   "Manage IAM groups and the privileges they grant",
			Long: `A Group is the central building block of IONOS CLOUD Identity & Access Management (IAM). In the IAM model there are three moving parts:

  Users     - individual cloud accounts (see ` + "`ionosctl compute user`" + `).
  Groups    - a named collection of Users that all share the same set of contract-wide PRIVILEGES.
  Resources - datacenters, snapshots, images, IP blocks, etc. that a Group can be granted access to (see ` + "`ionosctl compute share`" + `).

Privileges vs. shares - two different things:
  * Privileges (this command's boolean flags, e.g. --create-dc, --reserve-ip, --create-k8s) are CONTRACT-WIDE CAPABILITIES. They decide what kinds of actions the Group's members may perform anywhere on the contract (for example "may create datacenters").
  * A Share (` + "`ionosctl compute group resource`" + ` to list, ` + "`ionosctl compute share`" + ` to grant) gives the Group access to ONE SPECIFIC existing resource, with optional edit/re-share rights.

A User inherits the UNION of the privileges of every Group they belong to. A User flagged as administrator (--administrator on ` + "`ionosctl compute user`" + `) bypasses group privileges entirely and has full contract access.

The sub-commands here let you list/get/create/update/delete Groups, manage their membership (` + "`group user`" + `), and see the resources a Group can reach (` + "`group resource`" + `).`,
			TraverseChildren: true,
		},
	}
	groupCmd.AddColsFlag(allGroupCols)

	groupCmd.AddCommand(GroupListCmd())
	groupCmd.AddCommand(GroupGetCmd())
	groupCmd.AddCommand(GroupCreateCmd())
	groupCmd.AddCommand(GroupUpdateCmd())
	groupCmd.AddCommand(GroupDeleteCmd())
	groupCmd.AddCommand(GroupResourceCmd())
	groupCmd.AddCommand(GroupUserCmd())

	return core.WithConfigOverride(groupCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
