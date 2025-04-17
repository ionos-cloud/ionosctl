package customDomains

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"CustomDomainsId", "Name", "CertificateId"}
)

func CustomDomainsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "customdomains",
			Short:            "An array of custom domains",
			Aliases:          []string{"custom-domains"},
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(AddCmd())
	cmd.AddCommand(ListCmd())
	cmd.AddCommand(RemovetCmd())
	return cmd
}
