package share

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ShareCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "share",
		Resource:  "share",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Grant a Group access to a specific resource",
		LongDesc: `Grant a Group access to one specific existing resource, creating a Share for the (Group, Resource) pair. Use this to hand a Group a concrete datacenter, image, snapshot, IP block, etc. - separate from the contract-wide privileges you set on the Group itself.

By default the share grants read/use access only. Add --edit-privilege to let members modify the resource, and/or --share-privilege to let them re-share it with other Groups. Find shareable resource IDs with ` + "`ionosctl compute resource list`" + `.

Required values to run a command:

* Group Id
* Resource Id`,
		Example: `# Give a group read/use access to a datacenter
ionosctl compute share create --group-id GROUP_ID --resource-id DATACENTER_ID

# Give a group full control: edit the resource and re-share it
ionosctl compute share create --group-id GROUP_ID --resource-id RESOURCE_ID --edit-privilege --share-privilege`,
		PreCmdRun:  PreRunGroupResourceIds,
		CmdRun:     RunShareCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ResourcesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgEditPrivilege, "", false, "Also allow the Group's members to edit (modify) the shared resource, not just view/use it")
	cmd.AddBoolFlag(cloudapiv6.ArgSharePrivilege, "", false, "Also allow the Group's members to re-share this resource with other Groups")

	return cmd
}
