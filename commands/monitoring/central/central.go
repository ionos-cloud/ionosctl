package central

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func CentralCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "central",
			Aliases:          []string{"c"},
			Short:            "Central monitoring is a feature that allows you to use other products to send metrics to a central location.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(CentralFindByIdCmd())
	cmd.AddCommand(CentralDisable())
	cmd.AddCommand(CentralEnable())

	return cmd
}
