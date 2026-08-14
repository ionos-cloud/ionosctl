package server

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/cdrom"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/console"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/gpu"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/token"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/volume"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

const (
	serverCubeType       = "CUBE"
	serverEnterpriseType = "ENTERPRISE"
	serverVCPUType       = "VCPU"
	serverGPUType        = "GPU"
)

// AllServerCols defines the columns for server output. Exported for vm-autoscaling.
var AllServerCols = []table.Column{
	{Name: "ServerId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Type", JSONPath: "properties.type", Default: true},
	{Name: "AvailabilityZone", JSONPath: "properties.availabilityZone", Default: true},
	{Name: "Cores", JSONPath: "properties.cores", Default: true},
	{Name: "RAM", JSONPath: "properties.ram", Default: true},
	{Name: "CpuFamily", JSONPath: "properties.cpuFamily", Default: true},
	{Name: "VmState", JSONPath: "properties.vmState", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "DatacenterId", JSONPath: "href"},
	{Name: "TemplateId", JSONPath: "properties.templateUuid"},
	{Name: "BootCdromId", JSONPath: "properties.bootCdrom.id"},
	{Name: "BootVolumeId", JSONPath: "properties.bootVolume.id"},
	{Name: "NicMultiQueue", JSONPath: "properties.nicMultiQueue"},
	{Name: "EnabledFeatures", JSONPath: "properties.enabledFeatures"},
}

func ServerCmd() *core.Command {
	serverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "server",
			Aliases:          []string{"s", "svr"},
			Short:            "Server Operations",
			Long:             "The sub-commands of `ionosctl compute server` allow you to manage Servers.",
			TraverseChildren: true,
		},
	}
	serverCmd.AddColsFlag(AllServerCols)

	serverCmd.AddCommand(ServerListCmd())
	serverCmd.AddCommand(ServerGetCmd())
	serverCmd.AddCommand(ServerCreateCmd())
	serverCmd.AddCommand(ServerUpdateCmd())
	serverCmd.AddCommand(ServerDeleteCmd())
	serverCmd.AddCommand(ServerStartCmd())
	serverCmd.AddCommand(ServerStopCmd())
	serverCmd.AddCommand(ServerRebootCmd())
	serverCmd.AddCommand(ServerSuspendCmd())
	serverCmd.AddCommand(ServerResumeCmd())

	serverCmd.AddCommand(token.ServerTokenCmd())
	serverCmd.AddCommand(console.ServerConsoleCmd())
	serverCmd.AddCommand(volume.ServerVolumeCmd())
	serverCmd.AddCommand(cdrom.ServerCdromCmd())
	serverCmd.AddCommand(gpu.ServerGpuCmd())

	return core.WithConfigOverride(serverCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
