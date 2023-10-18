package autoscaling

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

func AutoscalingCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "autoscaling",
			Aliases:          []string{"as", "autoscale"},
			Short:            "VM Autoscaling Operations",
			TraverseChildren: true,
		},
	}
	return cmd
}
