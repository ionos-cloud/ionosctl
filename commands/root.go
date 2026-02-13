package commands

import (
	api_gateway "github.com/ionos-cloud/ionosctl/v6/commands/api-gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/image"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn"
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka"

	certificates "github.com/ionos-cloud/ionosctl/v6/commands/cert"
	"github.com/ionos-cloud/ionosctl/v6/commands/cfg"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	container_registry "github.com/ionos-cloud/ionosctl/v6/commands/container-registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas"
	_ "github.com/ionos-cloud/ionosctl/v6/commands/dns" // registers itself via init()
	logging_service "github.com/ionos-cloud/ionosctl/v6/commands/logging-service"
	"github.com/ionos-cloud/ionosctl/v6/commands/token"
	vm_autoscaling "github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn"
	"github.com/ionos-cloud/ionosctl/v6/internal/cmd"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd.Execute()
}

func GetRootCmd() *core.Command {
	return cmd.GetRootCmd()
}

func init() {
	// Add SubCommands to RootCmd
	addCommands()
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	rootCmd := cmd.GetRootCmd()

	rootCmd.AddCommand(Shell())
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(Man())
	rootCmd.AddCommand(cfg.Login())

	// cfg
	rootCmd.AddCommand(cfg.ConfigCmd())
	// Config namespace commands are also available via the root command, but are hidden
	for _, c := range cfg.ConfigCmd().SubCommands() {
		if c.Name() == "location" {
			// This one is confusing without `cfg` namespace;
			// It also would override CPU Architecture locations command, so skip it.
			continue
		}
		c.Command.Hidden = true
		rootCmd.AddCommand(c)
	}

	// V6 Resources Commands
	rootCmd.AddCommand(cloudapiv6.LocationCmd())
	rootCmd.AddCommand(cloudapiv6.DatacenterCmd())
	rootCmd.AddCommand(cloudapiv6.ServerCmd())
	rootCmd.AddCommand(cloudapiv6.VolumeCmd())
	rootCmd.AddCommand(cloudapiv6.LanCmd())
	rootCmd.AddCommand(cloudapiv6.NatgatewayCmd())
	rootCmd.AddCommand(cloudapiv6.ApplicationLoadBalancerCmd())
	rootCmd.AddCommand(cloudapiv6.NetworkloadbalancerCmd())
	rootCmd.AddCommand(cloudapiv6.NicCmd())
	rootCmd.AddCommand(cloudapiv6.LoadBalancerCmd())
	rootCmd.AddCommand(cloudapiv6.IpblockCmd())
	rootCmd.AddCommand(cloudapiv6.IpconsumerCmd())
	rootCmd.AddCommand(cloudapiv6.IpfailoverCmd())
	rootCmd.AddCommand(cloudapiv6.RequestCmd())
	rootCmd.AddCommand(cloudapiv6.SnapshotCmd())
	rootCmd.AddCommand(image.ImageCmd())
	rootCmd.AddCommand(cloudapiv6.FirewallruleCmd())
	rootCmd.AddCommand(cloudapiv6.FlowlogCmd())
	rootCmd.AddCommand(cloudapiv6.LabelCmd())
	rootCmd.AddCommand(cloudapiv6.ContractCmd())
	rootCmd.AddCommand(cloudapiv6.UserCmd())
	rootCmd.AddCommand(cloudapiv6.GroupCmd())
	rootCmd.AddCommand(cloudapiv6.ResourceCmd())
	rootCmd.AddCommand(cloudapiv6.BackupunitCmd())
	rootCmd.AddCommand(cloudapiv6.PccCmd())
	rootCmd.AddCommand(cloudapiv6.ShareCmd())
	rootCmd.AddCommand(cloudapiv6.K8sCmd())
	rootCmd.AddCommand(cloudapiv6.TargetGroupCmd())
	rootCmd.AddCommand(cloudapiv6.TemplateCmd())
	// Auth Command
	rootCmd.AddCommand(token.TokenCmd())
	// Add DBaaS Commands
	rootCmd.AddCommand(dbaas.DataBaseServiceCmd())
	// Add Certificate Manager Commands
	rootCmd.AddCommand(certificates.Root())
	// Add Container Registry Commands
	rootCmd.AddCommand(container_registry.ContainerRegistryCmd())
	// VM-Autoscaling commands
	rootCmd.AddCommand(vm_autoscaling.Root())

	// dns registers itself via init()
	rootCmd.AddCommand(logging_service.Root())

	rootCmd.AddCommand(api_gateway.Root())

	rootCmd.AddCommand(monitoring.Root())

	rootCmd.AddCommand(cdn.Command())

	rootCmd.AddCommand(vpn.Root())

	rootCmd.AddCommand(kafka.Command())
}
