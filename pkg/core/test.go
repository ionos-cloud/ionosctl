package core

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6/mocks"
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
		NS:        testConst,
		Namespace: testConst,
		Resource:  testConst,
		Verb:      testConst,
		Printer:   prt,
	}
	preRunner(preCmdCfg)
}

type CmdRunnerTest func(c *CommandConfig, mocks *ResourcesMocks)

type ResourcesMocks struct {
	Client              *mock_resources.MockClientService
	Location            *mock_resources.MockLocationsService
	Datacenter          *mock_resources.MockDatacentersService
	Server              *mock_resources.MockServersService
	Volume              *mock_resources.MockVolumesService
	Lan                 *mock_resources.MockLansService
	NatGateway          *mock_resources.MockNatGatewaysService
	NetworkLoadBalancer *mock_resources.MockNetworkLoadBalancersService
	Nic                 *mock_resources.MockNicsService
	Loadbalancer        *mock_resources.MockLoadbalancersService
	IpBlocks            *mock_resources.MockIpBlocksService
	Request             *mock_resources.MockRequestsService
	Image               *mock_resources.MockImagesService
	Snapshot            *mock_resources.MockSnapshotsService
	FirewallRule        *mock_resources.MockFirewallRulesService
	FlowLog             *mock_resources.MockFlowLogsService
	Label               *mock_resources.MockLabelResourcesService
	Contract            *mock_resources.MockContractsService
	User                *mock_resources.MockUsersService
	Group               *mock_resources.MockGroupsService
	S3Key               *mock_resources.MockS3KeysService
	BackupUnit          *mock_resources.MockBackupUnitsService
	Pcc                 *mock_resources.MockPccsService
	K8s                 *mock_resources.MockK8sService
	Template            *mock_resources.MockTemplatesService
}

func CmdConfigTest(t *testing.T, writer io.Writer, runner CmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	printReg, _ := printer.NewPrinterRegistry(writer, writer)
	prt := printReg[viper.GetString(config.ArgOutput)]
	// Init Test Mock Resources and Services
	testMocks := initMockResources(ctrl)
	cmdConfig := &CommandConfig{
		NS:        testConst,
		Namespace: testConst,
		Resource:  testConst,
		Verb:      testConst,
		Printer:   prt,
		Context:   context.TODO(),
		initCfg:   func(c *CommandConfig) error { return nil },
	}
	cmdConfig = initMockServices(cmdConfig, testMocks)
	runner(cmdConfig, testMocks)
}

// Init Mock Resources for Test
func initMockResources(ctrl *gomock.Controller) *ResourcesMocks {
	return &ResourcesMocks{
		Client:              mock_resources.NewMockClientService(ctrl),
		Location:            mock_resources.NewMockLocationsService(ctrl),
		Datacenter:          mock_resources.NewMockDatacentersService(ctrl),
		Server:              mock_resources.NewMockServersService(ctrl),
		Lan:                 mock_resources.NewMockLansService(ctrl),
		Volume:              mock_resources.NewMockVolumesService(ctrl),
		NatGateway:          mock_resources.NewMockNatGatewaysService(ctrl),
		NetworkLoadBalancer: mock_resources.NewMockNetworkLoadBalancersService(ctrl),
		Nic:                 mock_resources.NewMockNicsService(ctrl),
		Loadbalancer:        mock_resources.NewMockLoadbalancersService(ctrl),
		IpBlocks:            mock_resources.NewMockIpBlocksService(ctrl),
		Request:             mock_resources.NewMockRequestsService(ctrl),
		Image:               mock_resources.NewMockImagesService(ctrl),
		Snapshot:            mock_resources.NewMockSnapshotsService(ctrl),
		FirewallRule:        mock_resources.NewMockFirewallRulesService(ctrl),
		FlowLog:             mock_resources.NewMockFlowLogsService(ctrl),
		Label:               mock_resources.NewMockLabelResourcesService(ctrl),
		Contract:            mock_resources.NewMockContractsService(ctrl),
		User:                mock_resources.NewMockUsersService(ctrl),
		Group:               mock_resources.NewMockGroupsService(ctrl),
		S3Key:               mock_resources.NewMockS3KeysService(ctrl),
		BackupUnit:          mock_resources.NewMockBackupUnitsService(ctrl),
		Pcc:                 mock_resources.NewMockPccsService(ctrl),
		K8s:                 mock_resources.NewMockK8sService(ctrl),
		Template:            mock_resources.NewMockTemplatesService(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocks) *CommandConfig {
	c.Locations = func() v6.LocationsService { return tm.Location }
	c.DataCenters = func() v6.DatacentersService { return tm.Datacenter }
	c.Servers = func() v6.ServersService { return tm.Server }
	c.Volumes = func() v6.VolumesService { return tm.Volume }
	c.Lans = func() v6.LansService { return tm.Lan }
	c.NatGateways = func() v6.NatGatewaysService { return tm.NatGateway }
	c.NetworkLoadBalancers = func() v6.NetworkLoadBalancersService { return tm.NetworkLoadBalancer }
	c.Nics = func() v6.NicsService { return tm.Nic }
	c.Loadbalancers = func() v6.LoadbalancersService { return tm.Loadbalancer }
	c.IpBlocks = func() v6.IpBlocksService { return tm.IpBlocks }
	c.Requests = func() v6.RequestsService { return tm.Request }
	c.Images = func() v6.ImagesService { return tm.Image }
	c.Snapshots = func() v6.SnapshotsService { return tm.Snapshot }
	c.FirewallRules = func() v6.FirewallRulesService { return tm.FirewallRule }
	c.FlowLogs = func() v6.FlowLogsService { return tm.FlowLog }
	c.Labels = func() v6.LabelResourcesService { return tm.Label }
	c.Contracts = func() v6.ContractsService { return tm.Contract }
	c.Users = func() v6.UsersService { return tm.User }
	c.Groups = func() v6.GroupsService { return tm.Group }
	c.S3Keys = func() v6.S3KeysService { return tm.S3Key }
	c.BackupUnit = func() v6.BackupUnitsService { return tm.BackupUnit }
	c.Pccs = func() v6.PccsService { return tm.Pcc }
	c.K8s = func() v6.K8sService { return tm.K8s }
	c.Templates = func() v6.TemplatesService { return tm.Template }
	return c
}
