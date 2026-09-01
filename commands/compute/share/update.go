package share

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ShareUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "share",
		Resource:  "share",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Change a Group's permissions on a shared resource",
		LongDesc: `Change the permission bits (--edit-privilege, --share-privilege) of an existing Share for a (Group, Resource) pair. Use this to promote a read-only share to editable/re-shareable, or to walk those permissions back.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Group Id
* Resource Id`,
		Example: `# Allow the group to re-share the resource
ionosctl compute share update --group-id GROUP_ID --resource-id RESOURCE_ID --share-privilege

# Revoke edit rights but keep the share in place
ionosctl compute share update --group-id GROUP_ID --resource-id RESOURCE_ID --edit-privilege=false`,
		PreCmdRun:  PreRunGroupResourceIds,
		CmdRun:     RunShareUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cobraCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgEditPrivilege, "", false, "Set whether the Group's members may edit (modify) the shared resource. E.g.: --edit-privilege=true, --edit-privilege=false")
	cmd.AddBoolFlag(cloudapiv6.ArgSharePrivilege, "", false, "Set whether the Group's members may re-share this resource with other Groups. E.g.: --share-privilege=true, --share-privilege=false")

	return cmd
}
