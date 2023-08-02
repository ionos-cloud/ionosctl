package cloudapi_v6

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	mockResources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources/mocks"
)

type ResourcesMocks struct {
	Location                *mockResources.MockLocationsService
	Datacenter              *mockResources.MockDatacentersService
	Server                  *mockResources.MockServersService
	Volume                  *mockResources.MockVolumesService
	Lan                     *mockResources.MockLansService
	NatGateway              *mockResources.MockNatGatewaysService
	ApplicationLoadBalancer *mockResources.MockApplicationLoadBalancersService
	NetworkLoadBalancer     *mockResources.MockNetworkLoadBalancersService
	Nic                     *mockResources.MockNicsService
	Loadbalancer            *mockResources.MockLoadbalancersService
	IpBlocks                *mockResources.MockIpBlocksService
	Request                 *mockResources.MockRequestsService
	Image                   *mockResources.MockImagesService
	Snapshot                *mockResources.MockSnapshotsService
	FirewallRule            *mockResources.MockFirewallRulesService
	FlowLog                 *mockResources.MockFlowLogsService
	Label                   *mockResources.MockLabelResourcesService
	Contract                *mockResources.MockContractsService
	User                    *mockResources.MockUsersService
	Group                   *mockResources.MockGroupsService
	S3Key                   *mockResources.MockS3KeysService
	BackupUnit              *mockResources.MockBackupUnitsService
	Pcc                     *mockResources.MockPccsService
	K8s                     *mockResources.MockK8sService
	Template                *mockResources.MockTemplatesService
	TargetGroup             *mockResources.MockTargetGroupsService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Location:                mockResources.NewMockLocationsService(ctrl),
		Datacenter:              mockResources.NewMockDatacentersService(ctrl),
		Server:                  mockResources.NewMockServersService(ctrl),
		Lan:                     mockResources.NewMockLansService(ctrl),
		Volume:                  mockResources.NewMockVolumesService(ctrl),
		NatGateway:              mockResources.NewMockNatGatewaysService(ctrl),
		ApplicationLoadBalancer: mockResources.NewMockApplicationLoadBalancersService(ctrl),
		NetworkLoadBalancer:     mockResources.NewMockNetworkLoadBalancersService(ctrl),
		Nic:                     mockResources.NewMockNicsService(ctrl),
		Loadbalancer:            mockResources.NewMockLoadbalancersService(ctrl),
		IpBlocks:                mockResources.NewMockIpBlocksService(ctrl),
		Request:                 mockResources.NewMockRequestsService(ctrl),
		Image:                   mockResources.NewMockImagesService(ctrl),
		Snapshot:                mockResources.NewMockSnapshotsService(ctrl),
		FirewallRule:            mockResources.NewMockFirewallRulesService(ctrl),
		FlowLog:                 mockResources.NewMockFlowLogsService(ctrl),
		Label:                   mockResources.NewMockLabelResourcesService(ctrl),
		Contract:                mockResources.NewMockContractsService(ctrl),
		User:                    mockResources.NewMockUsersService(ctrl),
		Group:                   mockResources.NewMockGroupsService(ctrl),
		S3Key:                   mockResources.NewMockS3KeysService(ctrl),
		BackupUnit:              mockResources.NewMockBackupUnitsService(ctrl),
		Pcc:                     mockResources.NewMockPccsService(ctrl),
		K8s:                     mockResources.NewMockK8sService(ctrl),
		Template:                mockResources.NewMockTemplatesService(ctrl),
		TargetGroup:             mockResources.NewMockTargetGroupsService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Locations = func() resources.LocationsService { return tm.Location }
	c.DataCenters = func() resources.DatacentersService { return tm.Datacenter }
	c.Servers = func() resources.ServersService { return tm.Server }
	c.Volumes = func() resources.VolumesService { return tm.Volume }
	c.Lans = func() resources.LansService { return tm.Lan }
	c.NatGateways = func() resources.NatGatewaysService { return tm.NatGateway }
	c.ApplicationLoadBalancers = func() resources.ApplicationLoadBalancersService { return tm.ApplicationLoadBalancer }
	c.NetworkLoadBalancers = func() resources.NetworkLoadBalancersService { return tm.NetworkLoadBalancer }
	c.Nics = func() resources.NicsService { return tm.Nic }
	c.Loadbalancers = func() resources.LoadbalancersService { return tm.Loadbalancer }
	c.IpBlocks = func() resources.IpBlocksService { return tm.IpBlocks }
	c.Requests = func() resources.RequestsService { return tm.Request }
	c.Images = func() resources.ImagesService { return tm.Image }
	c.Snapshots = func() resources.SnapshotsService { return tm.Snapshot }
	c.FirewallRules = func() resources.FirewallRulesService { return tm.FirewallRule }
	c.FlowLogs = func() resources.FlowLogsService { return tm.FlowLog }
	c.Labels = func() resources.LabelResourcesService { return tm.Label }
	c.Contracts = func() resources.ContractsService { return tm.Contract }
	c.Users = func() resources.UsersService { return tm.User }
	c.Groups = func() resources.GroupsService { return tm.Group }
	c.S3Keys = func() resources.S3KeysService { return tm.S3Key }
	c.BackupUnit = func() resources.BackupUnitsService { return tm.BackupUnit }
	c.Pccs = func() resources.PccsService { return tm.Pcc }
	c.K8s = func() resources.K8sService { return tm.K8s }
	c.Templates = func() resources.TemplatesService { return tm.Template }
	c.TargetGroups = func() resources.TargetGroupsService { return tm.TargetGroup }
	return c
}
