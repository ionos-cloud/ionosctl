package datacenter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func DatacenterCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Virtual Data Center in a chosen region",
		LongDesc: `Create a Virtual Data Center (VDC) - the top-level, network-isolated container that will hold your compute resources (servers, volumes, LANs, NICs, firewalls).

Everything you provision afterwards is created inside a VDC and inherits its region, so the single most important choice here is ` + "`--location`" + `: the region the VDC lives in. This is set once at creation and CANNOT be changed later - to run workloads in another region you create a separate VDC there. You can provision as many VDCs as your contract allows; each is logically segmented from the others.

The name defaults to "Unnamed Data Center" and the location defaults to ` + "`de/txl`" + ` (Berlin), so a VDC can be created with no flags at all - but it is strongly recommended to pass an explicit ` + "`--location`" + ` so you are not surprised by where your resources land.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the VDC reaches the AVAILABLE state before the command returns; without it the command returns as soon as the provisioning request is accepted.`,
		Example: `# Create a VDC in Frankfurt with an explicit name
ionosctl compute datacenter create --name "prod-vdc" --location de/fra

# Create a VDC and wait until it is AVAILABLE, with a description, showing only chosen columns
ionosctl compute datacenter create --name "prod-vdc" --description "Production workloads, EU" --location de/fra --wait --cols "DatacenterId,Name,Location,State"`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunDataCenterCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Data Center", "Human-friendly name shown in the DCD and CLI listings. Does not need to be unique")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Free-text description of the VDC's purpose. Optional and editable later")
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "de/txl", "Region the VDC and all resources inside it will live in, e.g. de/txl (Berlin), de/fra (Frankfurt), gb/lhr (London), es/vit (Logrono), fr/par (Paris), us/las (Las Vegas), us/ewr (Newark). IMMUTABLE - cannot be changed after creation. Must be enabled for your contract")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
