package ipblock

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "IpBlockId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "Size", JSONPath: "properties.size", Default: true},
	{Name: "Ips", JSONPath: "properties.ips", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func IpblockCmd() *core.Command {
	ipblockCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipblock",
			Aliases:          []string{"ip", "ipb"},
			Short:            "IpBlock Operations",
			Long:             "The sub-commands of `ionosctl compute ipblock` allow you to create/reserve, list, get, update, delete IpBlocks.",
			TraverseChildren: true,
		},
	}
	ipblockCmd.AddColsFlag(allCols)

	ipblockCmd.AddCommand(IpBlockListCmd())
	ipblockCmd.AddCommand(IpBlockGetCmd())
	ipblockCmd.AddCommand(IpBlockCreateCmd())
	ipblockCmd.AddCommand(IpBlockUpdateCmd())
	ipblockCmd.AddCommand(IpBlockDeleteCmd())

	return core.WithConfigOverride(ipblockCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
