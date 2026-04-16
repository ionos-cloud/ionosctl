package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Short:            "The sub-commands of 'ionosctl kafka user' allow you to manage kafka users",
			Aliases:          []string{"u"},
			TraverseChildren: true,
		},
	}
	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(GetAccess())
	return cmd
}
