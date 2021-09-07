package cloudapi_v5

import (
	"context"

	config2 "github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/viper"
)

type Services struct {
	// Compute Resources Services
	Locations     func() resources.LocationsService
	DataCenters   func() resources.DatacentersService
	Servers       func() resources.ServersService
	Volumes       func() resources.VolumesService
	Lans          func() resources.LansService
	Nics          func() resources.NicsService
	Loadbalancers func() resources.LoadbalancersService
	Requests      func() resources.RequestsService
	Images        func() resources.ImagesService
	Snapshots     func() resources.SnapshotsService
	IpBlocks      func() resources.IpBlocksService
	FirewallRules func() resources.FirewallRulesService
	Labels        func() resources.LabelResourcesService
	Contracts     func() resources.ContractsService
	Users         func() resources.UsersService
	Groups        func() resources.GroupsService
	S3Keys        func() resources.S3KeysService
	BackupUnit    func() resources.BackupUnitsService
	Pccs          func() resources.PccsService
	K8s           func() resources.K8sService
	// Context
	Context context.Context
}

// InitClient for Commands
func (c *Services) InitClient() (*resources.Client, error) {
	clientSvc, err := resources.NewClientService(
		viper.GetString(config2.Username),
		viper.GetString(config2.Password),
		viper.GetString(config2.Token), // Token support
		config2.GetServerUrl(),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}

// InitServices for Commands
func (c *Services) InitServices(client *resources.Client) error {
	c.Locations = func() resources.LocationsService { return resources.NewLocationService(client, c.Context) }
	c.DataCenters = func() resources.DatacentersService { return resources.NewDataCenterService(client, c.Context) }
	c.Servers = func() resources.ServersService { return resources.NewServerService(client, c.Context) }
	c.Volumes = func() resources.VolumesService { return resources.NewVolumeService(client, c.Context) }
	c.Lans = func() resources.LansService { return resources.NewLanService(client, c.Context) }
	c.Nics = func() resources.NicsService { return resources.NewNicService(client, c.Context) }
	c.Loadbalancers = func() resources.LoadbalancersService { return resources.NewLoadbalancerService(client, c.Context) }
	c.IpBlocks = func() resources.IpBlocksService { return resources.NewIpBlockService(client, c.Context) }
	c.Requests = func() resources.RequestsService { return resources.NewRequestService(client, c.Context) }
	c.Images = func() resources.ImagesService { return resources.NewImageService(client, c.Context) }
	c.Snapshots = func() resources.SnapshotsService { return resources.NewSnapshotService(client, c.Context) }
	c.FirewallRules = func() resources.FirewallRulesService { return resources.NewFirewallRuleService(client, c.Context) }
	c.Labels = func() resources.LabelResourcesService { return resources.NewLabelResourceService(client, c.Context) }
	c.Contracts = func() resources.ContractsService { return resources.NewContractService(client, c.Context) }
	c.Users = func() resources.UsersService { return resources.NewUserService(client, c.Context) }
	c.Groups = func() resources.GroupsService { return resources.NewGroupService(client, c.Context) }
	c.S3Keys = func() resources.S3KeysService { return resources.NewS3KeyService(client, c.Context) }
	c.BackupUnit = func() resources.BackupUnitsService { return resources.NewBackupUnitService(client, c.Context) }
	c.Pccs = func() resources.PccsService { return resources.NewPrivateCrossConnectService(client, c.Context) }
	c.K8s = func() resources.K8sService { return resources.NewK8sService(client, c.Context) }
	return nil
}
