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
			Use:              "group",
			Aliases:          []string{"g"},
			Short:            "Group Operations",
			Long:             "The sub-commands of `ionosctl compute group` allow you to list, get, create, update, delete Groups, but also operations: add/remove/list/update User from the Group.",
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
