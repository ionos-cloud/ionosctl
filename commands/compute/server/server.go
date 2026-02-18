package server

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/cdrom"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/console"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/gpu"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/token"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server/volume"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	serverCubeType       = "CUBE"
	serverEnterpriseType = "ENTERPRISE"
	serverVCPUType       = "VCPU"
	serverGPUType        = "GPU"
)

var (
	defaultServerCols = []string{"ServerId", "Name", "Type", "AvailabilityZone", "Cores", "RAM", "CpuFamily", "VmState", "State"}
	allServerCols     = []string{"ServerId", "DatacenterId", "Name", "AvailabilityZone", "Cores", "RAM", "CpuFamily", "VmState", "State", "TemplateId", "Type", "BootCdromId", "BootVolumeId", "NicMultiQueue"}
)

// AllServerCols is the exported version of allServerCols for external consumers.
var AllServerCols = allServerCols

func ServerCmd() *core.Command {
	serverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "server",
			Aliases:          []string{"s", "svr"},
			Short:            "Server Operations",
			Long:             "The sub-commands of `ionosctl server` allow you to manage Servers.",
			TraverseChildren: true,
		},
	}
	globalFlags := serverCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultServerCols, tabheaders.ColsMessage(allServerCols))
	_ = viper.BindPFlag(core.GetFlagName(serverCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = serverCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allServerCols, cobra.ShellCompDirectiveNoFileComp
	})

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
