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
			Use:   "share",
			Short: "Share individual resources with a Group",
			Long: `A Share grants a Group access to ONE specific existing resource - a datacenter, snapshot, image, IP block, PCC, backup unit, etc. This is how you delegate a concrete resource, as opposed to the contract-wide capabilities set on the Group itself (` + "`ionosctl compute group`" + ` privilege flags).

Every share is scoped to a (Group, Resource) pair and carries two optional permission bits:
  * --edit-privilege  - members of the Group may modify the shared resource, not just view/use it.
  * --share-privilege - members of the Group may re-share the resource with other Groups.
With neither bit set, the Group gets read/use access only.

Once shared, the resource shows up under ` + "`ionosctl compute group resource list`" + ` for that Group. Discover shareable resources and their IDs with ` + "`ionosctl compute resource`" + `.`,
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
