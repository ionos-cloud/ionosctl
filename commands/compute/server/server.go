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
			Use:     "server",
			Aliases: []string{"s", "svr"},
			Short:   "Create and manage Compute Engine Servers (virtual machines)",
			Long: `The sub-commands of ` + "`ionosctl compute server`" + ` manage Servers: the virtual machines that run inside a Virtual Data Center (VDC, addressed by --datacenter-id).

A Server has a compute shape and a lifecycle:

  * Type (--type) selects the compute model:
      - ENTERPRISE: you size the machine yourself with explicit --cores and --ram, and pick a --cpu-family (INTEL_SKYLAKE, INTEL_XEON, AMD_OPTERON, ...). The most flexible type.
      - VCPU: like ENTERPRISE but with vCPUs and no --cpu-family (the platform picks it). Cost-effective; good for dev/test and general workloads.
      - CUBE: a fixed bundle chosen by --template-id (see ` + "`ionosctl compute template`" + `). Cores, RAM and an included NVMe Direct Attached Storage (DAS) boot volume all come from the template, so you do NOT pass --cores/--ram. CUBE is the only type that can be suspended/resumed.
      - GPU: a fixed --template-id bundle with an attached GPU; CPU family is assigned automatically (AMD_TURIN).

  * Lifecycle (start/stop/reboot, plus suspend/resume for CUBE): a Server can be running or deallocated. See the individual verbs.`,
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
