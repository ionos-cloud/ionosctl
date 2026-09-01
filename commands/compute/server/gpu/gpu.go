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
			Use:     "gpu",
			Aliases: []string{"gpus"},
			Short:   "Inspect the GPUs attached to a GPU Server",
			Long:    "The sub-commands of `ionosctl compute server gpu` list and get the GPU cards attached to a Server of --type GPU. GPUs are provisioned as part of the server's instance-size template (`server create --type GPU --template-id ...`), so these commands are read-only: they report each card's vendor, model and state. Only servers of type GPU carry GPUs.",

			TraverseChildren: true,
		},
	}

	serverGpuCmd.AddCommand(ServerGpuListCmd())
	serverGpuCmd.AddCommand(ServerGpuGetCmd())

	return core.WithConfigOverride(serverGpuCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
