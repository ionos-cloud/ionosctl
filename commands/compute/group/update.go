package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
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
		ShortDesc: "Update a Group's name or privileges",
		LongDesc: `Update an existing IAM Group's name and/or privileges. Any privilege flag you set is applied to the Group and therefore to ALL of its members; privilege flags you do NOT pass keep their current value, so you can toggle a single capability without listing the rest.

Setting a flag to false REVOKES that capability from every member of the Group (unless another Group they belong to still grants it).

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Group Id`,
		Example: `# Rename a group
ionosctl compute group update --group-id GROUP_ID --name "Platform Team"

# Grant an extra capability without touching the group's other privileges
ionosctl compute group update --group-id GROUP_ID --reserve-ip

# Revoke the ability to create datacenters and Kubernetes clusters
ionosctl compute group update --group-id GROUP_ID --create-dc=false --create-k8s=false`,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The new name for the Group")
	addGroupPrivilegeFlags(cmd)

	return cmd
}
