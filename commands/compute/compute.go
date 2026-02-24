package compute

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/backupunit"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/contract"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/datacenter"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/firewallrule"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/group"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/image"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/ipblock"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/ipconsumer"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/ipfailover"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/label"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/lan"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/loadbalancer"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/location"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/nic"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/pcc"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/request"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/resource"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/server"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/share"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/targetgroup"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/template"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/user"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/volume"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "compute",
			Short:            "Compute Engine operations (servers, datacenters, volumes, networking, etc.)",
			Long:             "The sub-commands of `ionosctl compute` allow you to manage Compute Engine resources such as Servers, Data Centers, Volumes, LANs, NICs, Load Balancers, and more.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(location.LocationCmd())
	cmd.AddCommand(datacenter.DatacenterCmd())
	cmd.AddCommand(server.ServerCmd())
	cmd.AddCommand(volume.VolumeCmd())
	cmd.AddCommand(lan.LanCmd())
	cmd.AddCommand(natgateway.NatgatewayCmd())
	cmd.AddCommand(applicationloadbalancer.ApplicationLoadBalancerCmd())
	cmd.AddCommand(networkloadbalancer.NetworkloadbalancerCmd())
	cmd.AddCommand(nic.NicCmd())
	cmd.AddCommand(loadbalancer.LoadBalancerCmd())
	cmd.AddCommand(ipblock.IpblockCmd())
	cmd.AddCommand(ipconsumer.IpconsumerCmd())
	cmd.AddCommand(ipfailover.IpfailoverCmd())
	cmd.AddCommand(request.RequestCmd())
	cmd.AddCommand(snapshot.SnapshotCmd())
	cmd.AddCommand(image.ImageCmd())
	cmd.AddCommand(firewallrule.FirewallruleCmd())
	cmd.AddCommand(flowlog.FlowlogCmd())
	cmd.AddCommand(label.LabelCmd())
	cmd.AddCommand(contract.ContractCmd())
	cmd.AddCommand(user.UserCmd())
	cmd.AddCommand(group.GroupCmd())
	cmd.AddCommand(resource.ResourceCmd())
	cmd.AddCommand(backupunit.BackupunitCmd())
	cmd.AddCommand(pcc.PccCmd())
	cmd.AddCommand(share.ShareCmd())
	cmd.AddCommand(k8s.K8sCmd())
	cmd.AddCommand(targetgroup.TargetGroupCmd())
	cmd.AddCommand(template.TemplateCmd())

	return cmd
}

// HiddenAliases returns all compute commands as hidden commands
// for backward compatibility at root level.
func HiddenAliases() []*core.Command {
	cmds := []*core.Command{
		location.LocationCmd(),
		datacenter.DatacenterCmd(),
		server.ServerCmd(),
		volume.VolumeCmd(),
		lan.LanCmd(),
		natgateway.NatgatewayCmd(),
		applicationloadbalancer.ApplicationLoadBalancerCmd(),
		networkloadbalancer.NetworkloadbalancerCmd(),
		nic.NicCmd(),
		loadbalancer.LoadBalancerCmd(),
		ipblock.IpblockCmd(),
		ipconsumer.IpconsumerCmd(),
		ipfailover.IpfailoverCmd(),
		request.RequestCmd(),
		snapshot.SnapshotCmd(),
		image.ImageCmd(),
		firewallrule.FirewallruleCmd(),
		flowlog.FlowlogCmd(),
		label.LabelCmd(),
		contract.ContractCmd(),
		user.UserCmd(),
		group.GroupCmd(),
		resource.ResourceCmd(),
		backupunit.BackupunitCmd(),
		pcc.PccCmd(),
		share.ShareCmd(),
		k8s.K8sCmd(),
		targetgroup.TargetGroupCmd(),
		template.TemplateCmd(),
	}

	for _, cmd := range cmds {
		cmd.Command.Hidden = true
	}

	return cmds
}

