package cloudapi_v5

import (
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	mockResources "github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources/mocks"
)

type ResourcesMocks struct {
	Client       *mockResources.MockClientService
	Location     *mockResources.MockLocationsService
	Datacenter   *mockResources.MockDatacentersService
	Server       *mockResources.MockServersService
	Volume       *mockResources.MockVolumesService
	Lan          *mockResources.MockLansService
	Nic          *mockResources.MockNicsService
	Loadbalancer *mockResources.MockLoadbalancersService
	IpBlocks     *mockResources.MockIpBlocksService
	Request      *mockResources.MockRequestsService
	Image        *mockResources.MockImagesService
	Snapshot     *mockResources.MockSnapshotsService
	FirewallRule *mockResources.MockFirewallRulesService
	Label        *mockResources.MockLabelResourcesService
	Contract     *mockResources.MockContractsService
	User         *mockResources.MockUsersService
	Group        *mockResources.MockGroupsService
	S3Key        *mockResources.MockS3KeysService
	BackupUnit   *mockResources.MockBackupUnitsService
	Pcc          *mockResources.MockPccsService
	K8s          *mockResources.MockK8sService
}

// InitMocksResources for Test
func InitMocksResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Client:       mockResources.NewMockClientService(ctrl),
		Location:     mockResources.NewMockLocationsService(ctrl),
		Datacenter:   mockResources.NewMockDatacentersService(ctrl),
		Server:       mockResources.NewMockServersService(ctrl),
		Lan:          mockResources.NewMockLansService(ctrl),
		Volume:       mockResources.NewMockVolumesService(ctrl),
		Nic:          mockResources.NewMockNicsService(ctrl),
		Loadbalancer: mockResources.NewMockLoadbalancersService(ctrl),
		IpBlocks:     mockResources.NewMockIpBlocksService(ctrl),
		Request:      mockResources.NewMockRequestsService(ctrl),
		Image:        mockResources.NewMockImagesService(ctrl),
		Snapshot:     mockResources.NewMockSnapshotsService(ctrl),
		FirewallRule: mockResources.NewMockFirewallRulesService(ctrl),
		Label:        mockResources.NewMockLabelResourcesService(ctrl),
		Contract:     mockResources.NewMockContractsService(ctrl),
		User:         mockResources.NewMockUsersService(ctrl),
		Group:        mockResources.NewMockGroupsService(ctrl),
		S3Key:        mockResources.NewMockS3KeysService(ctrl),
		BackupUnit:   mockResources.NewMockBackupUnitsService(ctrl),
		Pcc:          mockResources.NewMockPccsService(ctrl),
		K8s:          mockResources.NewMockK8sService(ctrl),
	}
}

// InitMockServices for Command Test
func InitMockServices(c *Services, tm *ResourcesMocks) *Services {
	c.Locations = func() resources.LocationsService { return tm.Location }
	c.DataCenters = func() resources.DatacentersService { return tm.Datacenter }
	c.Servers = func() resources.ServersService { return tm.Server }
	c.Volumes = func() resources.VolumesService { return tm.Volume }
	c.Lans = func() resources.LansService { return tm.Lan }
	c.Nics = func() resources.NicsService { return tm.Nic }
	c.Loadbalancers = func() resources.LoadbalancersService { return tm.Loadbalancer }
	c.IpBlocks = func() resources.IpBlocksService { return tm.IpBlocks }
	c.Requests = func() resources.RequestsService { return tm.Request }
	c.Images = func() resources.ImagesService { return tm.Image }
	c.Snapshots = func() resources.SnapshotsService { return tm.Snapshot }
	c.FirewallRules = func() resources.FirewallRulesService { return tm.FirewallRule }
	c.Labels = func() resources.LabelResourcesService { return tm.Label }
	c.Contracts = func() resources.ContractsService { return tm.Contract }
	c.Users = func() resources.UsersService { return tm.User }
	c.Groups = func() resources.GroupsService { return tm.Group }
	c.S3Keys = func() resources.S3KeysService { return tm.S3Key }
	c.BackupUnit = func() resources.BackupUnitsService { return tm.BackupUnit }
	c.Pccs = func() resources.PccsService { return tm.Pcc }
	c.K8s = func() resources.K8sService { return tm.K8s }
	return c
}
