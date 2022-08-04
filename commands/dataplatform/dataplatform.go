package dataplatform

import (
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func DataPlatformServiceCmd() *core.Command {
	dataPlatformCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dataplatform",
			Aliases:          []string{"dp"},
			Short:            "Data Platform Operations",
			Long:             "The sub-commands of `ionosctl dataplatform` allow you to perform operations on Data Platform resources.",
			TraverseChildren: true,
		},
	}
	dataPlatformCmd.AddCommand(ClusterCmd())
	return dataPlatformCmd
}
