package builder

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	mocks "github.com/ionos-cloud/ionosctl/pkg/resources/mocks"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
)

type PreCmdRunnerTest func(c *PreCommandConfig)

func PreCmdConfigTest(t *testing.T, writer io.Writer, preRunner PreCmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p, _ := printer.NewPrinterRegistry(writer, writer)
	prt := p[viper.GetString(config.ArgOutput)]
	preCmdConfig := &PreCommandConfig{
		Name:       "test",
		ParentName: "test",
		Printer:    prt,
	}
	preRunner(preCmdConfig)
}

type CmdRunnerTest func(c *CommandConfig, mocks *ResourcesMocks)

type ResourcesMocks struct {
	Client       *mocks.MockClientService
	Location     *mocks.MockLocationsService
	Datacenter   *mocks.MockDatacentersService
	Server       *mocks.MockServersService
	Volume       *mocks.MockVolumesService
	Lan          *mocks.MockLansService
	Nic          *mocks.MockNicsService
	Loadbalancer *mocks.MockLoadbalancersService
	IpBlocks     *mocks.MockIpBlocksService
	Request      *mocks.MockRequestsService
	Image        *mocks.MockImagesService
	Snapshot     *mocks.MockSnapshotsService
	FirewallRule *mocks.MockFirewallRulesService
	Label        *mocks.MockLabelResourcesService
	Contract     *mocks.MockContractsService
	User         *mocks.MockUsersService
}

func CmdConfigTest(t *testing.T, writer io.Writer, runner CmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p, _ := printer.NewPrinterRegistry(writer, writer)
	prt := p[viper.GetString(config.ArgOutput)]

	tm := &ResourcesMocks{
		Client:       mocks.NewMockClientService(ctrl),
		Location:     mocks.NewMockLocationsService(ctrl),
		Datacenter:   mocks.NewMockDatacentersService(ctrl),
		Server:       mocks.NewMockServersService(ctrl),
		Lan:          mocks.NewMockLansService(ctrl),
		Volume:       mocks.NewMockVolumesService(ctrl),
		Nic:          mocks.NewMockNicsService(ctrl),
		Loadbalancer: mocks.NewMockLoadbalancersService(ctrl),
		IpBlocks:     mocks.NewMockIpBlocksService(ctrl),
		Request:      mocks.NewMockRequestsService(ctrl),
		Image:        mocks.NewMockImagesService(ctrl),
		Snapshot:     mocks.NewMockSnapshotsService(ctrl),
		FirewallRule: mocks.NewMockFirewallRulesService(ctrl),
		Label:        mocks.NewMockLabelResourcesService(ctrl),
		Contract:     mocks.NewMockContractsService(ctrl),
		User:         mocks.NewMockUsersService(ctrl),
	}

	cmdConfig := &CommandConfig{
		Name:    "test",
		Printer: prt,
		Context: context.TODO(),
		initSvc: func(c *CommandConfig) error { return nil },
		Locations: func() resources.LocationsService {
			return tm.Location
		},
		DataCenters: func() resources.DatacentersService {
			return tm.Datacenter
		},
		Servers: func() resources.ServersService {
			return tm.Server
		},
		Volumes: func() resources.VolumesService {
			return tm.Volume
		},
		Lans: func() resources.LansService {
			return tm.Lan
		},
		Nics: func() resources.NicsService {
			return tm.Nic
		},
		Loadbalancers: func() resources.LoadbalancersService {
			return tm.Loadbalancer
		},
		IpBlocks: func() resources.IpBlocksService {
			return tm.IpBlocks
		},
		Requests: func() resources.RequestsService {
			return tm.Request
		},
		Images: func() resources.ImagesService {
			return tm.Image
		},
		Snapshots: func() resources.SnapshotsService {
			return tm.Snapshot
		},
		FirewallRules: func() resources.FirewallRulesService {
			return tm.FirewallRule
		},
		Labels: func() resources.LabelResourcesService {
			return tm.Label
		},
		Contracts: func() resources.ContractsService {
			return tm.Contract
		},
		Users: func() resources.UsersService {
			return tm.User
		},
	}

	runner(cmdConfig, tm)
}
