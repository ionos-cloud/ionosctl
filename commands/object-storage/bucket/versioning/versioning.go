package versioning

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Name", JSONPath: "Name", Default: true},
	{Name: "Versioning", JSONPath: "Versioning", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "versioning",
			Aliases:          []string{"ver"},
			Short:            "Manage bucket versioning configuration",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(GetCmd())

	return cmd
}
