package lan

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allLanCols = []table.Column{
	{Name: "LanId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Public", JSONPath: "properties.public", Default: true},
	{Name: "PccId", JSONPath: "properties.pcc", Default: true},
	{Name: "IPv6CidrBlock", JSONPath: "properties.ipv6CidrBlock", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "DatacenterId", JSONPath: "href"},
}

func LanCmd() *core.Command {
	lanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Aliases:          []string{"l"},
			Short:            "LAN Operations",
			Long:             "The sub-commands of `ionosctl compute lan` allow you to create, list, get, update, delete LANs.",
			TraverseChildren: true,
		},
	}
	lanCmd.AddColsFlag(allLanCols)

	lanCmd.AddCommand(LanListCmd())
	lanCmd.AddCommand(LanGetCmd())
	lanCmd.AddCommand(LanCreateCmd())
	lanCmd.AddCommand(LanUpdateCmd())
	lanCmd.AddCommand(LanDeleteCmd())

	return core.WithConfigOverride(lanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
