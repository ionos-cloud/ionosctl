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
			Use:     "pcc",
			Aliases: []string{"cc"},
			Short:   "Cross-Connect (Private Cross-Connect) Operations",
			Long: `The sub-commands of ` + "`" + `ionosctl compute pcc` + "`" + ` let you manage Cross-Connects (also called Private Cross-Connects, PCCs).

A Cross-Connect privately bridges LANs that live in different Virtual Data Centers so that servers on those LANs can talk to each other over the private network, without traffic ever leaving IONOS onto the public internet. Think of the Cross-Connect as the switch that peers the LANs together.

Key rules:
  * All connected LANs must belong to the same contract and the same region/location.
  * Only PRIVATE LANs can be bridged (a public LAN cannot join a Cross-Connect).
  * The private IP ranges of the connected LANs must not overlap (each LAN keeps its own addressing).
  * A given LAN can be attached to at most one Cross-Connect at a time.

This resource manages the Cross-Connect object itself (name, description) and lets you view which LANs are currently peered through it (` + "`" + `ionosctl compute pcc peers list` + "`" + `). Creating a Cross-Connect does NOT by itself connect any LAN. To actually attach a LAN, set that LAN's Cross-Connect with ` + "`" + `ionosctl compute lan update --pcc-id <PCC_ID>` + "`" + `.`,
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
