package upstreams

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"UpstreamId", "Scheme", "Loadbalancer", "Host", "Port", "Weight"}
)

func UpstreamsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "upstreams",
			Short:            "Manage route upstreams",
			Aliases:          []string{"streams"},
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(AddCmd())
	cmd.AddCommand(ListCmd())
	cmd.AddCommand(RemovetCmd())
	return cmd
}
