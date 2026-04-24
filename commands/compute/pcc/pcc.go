package pcc

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allPccCols = []table.Column{
	{Name: "PccId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func PccCmd() *core.Command {
	pccCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "pcc",
			Aliases:          []string{"cc"},
			Short:            "Cross-Connect Operations",
			Long:             "The sub-commands of `ionosctl compute pcc` allow you to list, get, create, update, delete Cross-Connect. To add Cross-Connect to a Lan, check the `ionosctl compute lan update` command.",
			TraverseChildren: true,
		},
	}
	pccCmd.AddColsFlag(allPccCols)

	pccCmd.AddCommand(PccListCmd())
	pccCmd.AddCommand(PccGetCmd())
	pccCmd.AddCommand(PccCreateCmd())
	pccCmd.AddCommand(PccUpdateCmd())
	pccCmd.AddCommand(PccDeleteCmd())
	pccCmd.AddCommand(PeersCmd())

	return core.WithConfigOverride(pccCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
