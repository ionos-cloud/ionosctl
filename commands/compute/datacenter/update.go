package datacenter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func DatacenterUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Virtual Data Center's name or description",
		LongDesc: `Update the editable properties of an existing Virtual Data Center: its ` + "`--name`" + ` and ` + "`--description`" + `.

Only these two fields can be changed. The VDC's region (` + "`location`" + `) is fixed at creation and is rejected by the API in update requests - to move workloads to another region you must create a new VDC there and recreate the resources. Renaming a VDC does not touch the resources inside it.

Pass only the flags you want to change; unspecified fields are left untouched.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the VDC is back in the AVAILABLE state.

Required values to run command:

* Data Center Id`,
		Example: `# Rename a VDC
ionosctl compute datacenter update --datacenter-id DATACENTER_ID --name "eu-prod"

# Change only the description and show the result
ionosctl compute datacenter update --datacenter-id DATACENTER_ID --description "Production workloads, EU" --cols "DatacenterId,Description"`,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunDataCenterUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIdShort, "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "New human-friendly name for the VDC. Omit to leave unchanged")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "New free-text description for the VDC. Omit to leave unchanged")

	return cmd
}
