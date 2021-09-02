package core

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	mockResources "github.com/ionos-cloud/ionosctl/pkg/resources/v6/mocks"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
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
		Command: &Command{
			Command: &cobra.Command{
				Use: testConst,
			},
		},
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
	Client              *mockResources.MockClientService
	Location            *mockResources.MockLocationsService
	Datacenter          *mockResources.MockDatacentersService
	Server              *mockResources.MockServersService
	Volume              *mockResources.MockVolumesService
	Lan                 *mockResources.MockLansService
	NatGateway          *mockResources.MockNatGatewaysService
	NetworkLoadBalancer *mockResources.MockNetworkLoadBalancersService
	Nic                 *mockResources.MockNicsService
	Loadbalancer        *mockResources.MockLoadbalancersService
	IpBlocks            *mockResources.MockIpBlocksService
	Request             *mockResources.MockRequestsService
	Image               *mockResources.MockImagesService
	Snapshot            *mockResources.MockSnapshotsService
	FirewallRule        *mockResources.MockFirewallRulesService
	FlowLog             *mockResources.MockFlowLogsService
	Label               *mockResources.MockLabelResourcesService
	Contract            *mockResources.MockContractsService
	User                *mockResources.MockUsersService
	Group               *mockResources.MockGroupsService
	S3Key               *mockResources.MockS3KeysService
	BackupUnit          *mockResources.MockBackupUnitsService
	Pcc                 *mockResources.MockPccsService
	K8s                 *mockResources.MockK8sService
	Template            *mockResources.MockTemplatesService
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
		Client:              mockResources.NewMockClientService(ctrl),
		Location:            mockResources.NewMockLocationsService(ctrl),
		Datacenter:          mockResources.NewMockDatacentersService(ctrl),
		Server:              mockResources.NewMockServersService(ctrl),
		Lan:                 mockResources.NewMockLansService(ctrl),
		Volume:              mockResources.NewMockVolumesService(ctrl),
		NatGateway:          mockResources.NewMockNatGatewaysService(ctrl),
		NetworkLoadBalancer: mockResources.NewMockNetworkLoadBalancersService(ctrl),
		Nic:                 mockResources.NewMockNicsService(ctrl),
		Loadbalancer:        mockResources.NewMockLoadbalancersService(ctrl),
		IpBlocks:            mockResources.NewMockIpBlocksService(ctrl),
		Request:             mockResources.NewMockRequestsService(ctrl),
		Image:               mockResources.NewMockImagesService(ctrl),
		Snapshot:            mockResources.NewMockSnapshotsService(ctrl),
		FirewallRule:        mockResources.NewMockFirewallRulesService(ctrl),
		FlowLog:             mockResources.NewMockFlowLogsService(ctrl),
		Label:               mockResources.NewMockLabelResourcesService(ctrl),
		Contract:            mockResources.NewMockContractsService(ctrl),
		User:                mockResources.NewMockUsersService(ctrl),
		Group:               mockResources.NewMockGroupsService(ctrl),
		S3Key:               mockResources.NewMockS3KeysService(ctrl),
		BackupUnit:          mockResources.NewMockBackupUnitsService(ctrl),
		Pcc:                 mockResources.NewMockPccsService(ctrl),
		K8s:                 mockResources.NewMockK8sService(ctrl),
		Template:            mockResources.NewMockTemplatesService(ctrl),
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
