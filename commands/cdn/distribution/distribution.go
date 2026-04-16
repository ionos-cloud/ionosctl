package distribution

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/cdn/distribution/routingrules"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Domain", JSONPath: "properties.domain", Default: true},
	{Name: "CertificateId", JSONPath: "properties.certificateId", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "distribution",
			Short:            "The sub-commands of 'ionosctl cdn distribution' allow you to manage CDN distributions",
			Aliases:          []string{"ds"},
			TraverseChildren: true,
		},
	}
	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(FindByID())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Create())
	cmd.AddCommand(Update())
	cmd.AddCommand(routingrules.Root())
	return cmd
}
