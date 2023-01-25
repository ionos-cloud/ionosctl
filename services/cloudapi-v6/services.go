package cloudapi_v6

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
)

type Services struct {
	// Compute Resources Services
	Locations                func() resources.LocationsService
	DataCenters              func() resources.DatacentersService
	Servers                  func() resources.ServersService
	Volumes                  func() resources.VolumesService
	Lans                     func() resources.LansService
	NatGateways              func() resources.NatGatewaysService
	ApplicationLoadBalancers func() resources.ApplicationLoadBalancersService
	NetworkLoadBalancers     func() resources.NetworkLoadBalancersService
	Nics                     func() resources.NicsService
	Loadbalancers            func() resources.LoadbalancersService
	Requests                 func() resources.RequestsService
	Images                   func() resources.ImagesService
	Snapshots                func() resources.SnapshotsService
	IpBlocks                 func() resources.IpBlocksService
	FirewallRules            func() resources.FirewallRulesService
	FlowLogs                 func() resources.FlowLogsService
	Labels                   func() resources.LabelResourcesService
	Contracts                func() resources.ContractsService
	Users                    func() resources.UsersService
	Groups                   func() resources.GroupsService
	S3Keys                   func() resources.S3KeysService
	BackupUnit               func() resources.BackupUnitsService
	Pccs                     func() resources.PccsService
	K8s                      func() resources.K8sService
	Templates                func() resources.TemplatesService
	TargetGroups             func() resources.TargetGroupsService
	// Context
	Context context.Context
}

// InitServices for Commands
func (c *Services) InitServices() error {
	client, err := config.GetClient()
	if err != nil {
		return err
	}
	c.Locations = func() resources.LocationsService { return resources.NewLocationService(client, c.Context) }
	c.DataCenters = func() resources.DatacentersService {
		return resources.NewDataCenterService(client, c.Context)
	}
	c.Servers = func() resources.ServersService { return resources.NewServerService(client, c.Context) }
	c.Volumes = func() resources.VolumesService { return resources.NewVolumeService(client, c.Context) }
	c.Lans = func() resources.LansService { return resources.NewLanService(client, c.Context) }
	c.NatGateways = func() resources.NatGatewaysService { return resources.NewNatGatewayService(client, c.Context) }
	c.ApplicationLoadBalancers = func() resources.ApplicationLoadBalancersService {
		return resources.NewApplicationLoadBalancerService(client, c.Context)
	}
	c.NetworkLoadBalancers = func() resources.NetworkLoadBalancersService {
		return resources.NewNetworkLoadBalancerService(client, c.Context)
	}
	c.Nics = func() resources.NicsService { return resources.NewNicService(client, c.Context) }
	c.Loadbalancers = func() resources.LoadbalancersService { return resources.NewLoadbalancerService(client, c.Context) }
	c.IpBlocks = func() resources.IpBlocksService { return resources.NewIpBlockService(client, c.Context) }
	c.Requests = func() resources.RequestsService { return resources.NewRequestService(client, c.Context) }
	c.Images = func() resources.ImagesService { return resources.NewImageService(client, c.Context) }
	c.Snapshots = func() resources.SnapshotsService { return resources.NewSnapshotService(client, c.Context) }
	c.FirewallRules = func() resources.FirewallRulesService { return resources.NewFirewallRuleService(client, c.Context) }
	c.FlowLogs = func() resources.FlowLogsService { return resources.NewFlowLogService(client, c.Context) }
	c.Labels = func() resources.LabelResourcesService { return resources.NewLabelResourceService(client, c.Context) }
	c.Contracts = func() resources.ContractsService { return resources.NewContractService(client, c.Context) }
	c.Users = func() resources.UsersService { return resources.NewUserService(client, c.Context) }
	c.Groups = func() resources.GroupsService { return resources.NewGroupService(client, c.Context) }
	c.S3Keys = func() resources.S3KeysService { return resources.NewS3KeyService(client, c.Context) }
	c.BackupUnit = func() resources.BackupUnitsService { return resources.NewBackupUnitService(client, c.Context) }
	c.Pccs = func() resources.PccsService { return resources.NewPrivateCrossConnectService(client, c.Context) }
	c.K8s = func() resources.K8sService { return resources.NewK8sService(client, c.Context) }
	c.Templates = func() resources.TemplatesService { return resources.NewTemplateService(client, c.Context) }
	c.TargetGroups = func() resources.TargetGroupsService {
		return resources.NewTargetGroupService(client, c.Context)
	}
	return nil
}
