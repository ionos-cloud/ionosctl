package gpu

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allGpuCols = []table.Column{
	{Name: "GpuId", JSONPath: "id", Default: true},
	{Name: "Type", JSONPath: "properties.type", Default: true},
	{Name: "Vendor", JSONPath: "properties.vendor", Default: true},
	{Name: "Model", JSONPath: "properties.model", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ServerGpuCmd() *core.Command {
	serverGpuCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gpu",
			Aliases:          []string{"gpus"},
			Short:            "GPU operations",
			Long:             "The sub-commands of `ionosctl compute server gpu` allow you to get and list Gpus from a Server.",
			TraverseChildren: true,
		},
	}

	serverGpuCmd.AddCommand(ServerGpuListCmd())
	serverGpuCmd.AddCommand(ServerGpuGetCmd())

	return core.WithConfigOverride(serverGpuCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
