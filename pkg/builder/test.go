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

const testConst = "test"

type PreCmdRunTest func(c *PreCommandConfig)

func PreCmdConfigTest(t *testing.T, writer io.Writer, preRunner PreCmdRunTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p, _ := printer.NewPrinterRegistry(writer, writer)
	prt := p[viper.GetString(config.ArgOutput)]
	preCmdCfg := &PreCommandConfig{
		Name:       testConst,
		ParentName: testConst,
		Printer:    prt,
	}
	preRunner(preCmdCfg)
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
	Group        *mocks.MockGroupsService
	S3Key        *mocks.MockS3KeysService
	Pcc          *mocks.MockPccsService
}

func CmdConfigTest(t *testing.T, writer io.Writer, runner CmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	printReg, _ := printer.NewPrinterRegistry(writer, writer)
	prt := printReg[viper.GetString(config.ArgOutput)]
	// Init Test Mock Resources and Services
	testMocks := initMockResources(ctrl)
	cmdConfig := &CommandConfig{
		Name:    testConst,
		Printer: prt,
		Context: context.TODO(),
		initCfg: func(c *CommandConfig) error { return nil },
	}
	cmdConfig = initMockServices(cmdConfig, testMocks)
	runner(cmdConfig, testMocks)
}

// Init Mock Resources for Test
func initMockResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
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
		Group:        mocks.NewMockGroupsService(ctrl),
		S3Key:        mocks.NewMockS3KeysService(ctrl),
		Pcc:          mocks.NewMockPccsService(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocks) *CommandConfig {
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
	c.Pccs = func() resources.PccsService { return tm.Pcc }
	return c
}
