package share

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allGroupShareCols = []table.Column{
	{Name: "ShareId", JSONPath: "id", Default: true},
	{Name: "EditPrivilege", JSONPath: "properties.editPrivilege", Default: true},
	{Name: "SharePrivilege", JSONPath: "properties.sharePrivilege", Default: true},
	{Name: "Type", JSONPath: "type", Default: true},
	{Name: "GroupId", JSONPath: "href"},
}

func ShareCmd() *core.Command {
	shareCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "share",
			Short:            "Resource Share Operations",
			Long:             "The sub-commands of `ionosctl compute share` allow you to list, get, create, update, delete Resource Shares.",
			TraverseChildren: true,
		},
	}
	shareCmd.AddColsFlag(allGroupShareCols)

	shareCmd.AddCommand(ShareListCmd())
	shareCmd.AddCommand(ShareGetCmd())
	shareCmd.AddCommand(ShareCreateCmd())
	shareCmd.AddCommand(ShareUpdateCmd())
	shareCmd.AddCommand(ShareDeleteCmd())

	return core.WithConfigOverride(shareCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
