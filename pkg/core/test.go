package core

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	mockResourcesv6 "github.com/ionos-cloud/ionosctl/pkg/resources/v6/mocks"
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
	Client       *mockResourcesv6.MockClientService
	Location     *mockResourcesv6.MockLocationsService
	Datacenter   *mockResourcesv6.MockDatacentersService
	Server       *mockResourcesv6.MockServersService
	Volume       *mockResourcesv6.MockVolumesService
	Lan          *mockResourcesv6.MockLansService
	Nic          *mockResourcesv6.MockNicsService
	Loadbalancer *mockResourcesv6.MockLoadbalancersService
	IpBlocks     *mockResourcesv6.MockIpBlocksService
	Request      *mockResourcesv6.MockRequestsService
	Image        *mockResourcesv6.MockImagesService
	Snapshot     *mockResourcesv6.MockSnapshotsService
	FirewallRule *mockResourcesv6.MockFirewallRulesService
	Label        *mockResourcesv6.MockLabelResourcesService
	Contract     *mockResourcesv6.MockContractsService
	User         *mockResourcesv6.MockUsersService
	Group        *mockResourcesv6.MockGroupsService
	S3Key        *mockResourcesv6.MockS3KeysService
	BackupUnit   *mockResourcesv6.MockBackupUnitsService
	Pcc          *mockResourcesv6.MockPccsService
	K8s          *mockResourcesv6.MockK8sService
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
		Client:       mockResourcesv6.NewMockClientService(ctrl),
		Location:     mockResourcesv6.NewMockLocationsService(ctrl),
		Datacenter:   mockResourcesv6.NewMockDatacentersService(ctrl),
		Server:       mockResourcesv6.NewMockServersService(ctrl),
		Lan:          mockResourcesv6.NewMockLansService(ctrl),
		Volume:       mockResourcesv6.NewMockVolumesService(ctrl),
		Nic:          mockResourcesv6.NewMockNicsService(ctrl),
		Loadbalancer: mockResourcesv6.NewMockLoadbalancersService(ctrl),
		IpBlocks:     mockResourcesv6.NewMockIpBlocksService(ctrl),
		Request:      mockResourcesv6.NewMockRequestsService(ctrl),
		Image:        mockResourcesv6.NewMockImagesService(ctrl),
		Snapshot:     mockResourcesv6.NewMockSnapshotsService(ctrl),
		FirewallRule: mockResourcesv6.NewMockFirewallRulesService(ctrl),
		Label:        mockResourcesv6.NewMockLabelResourcesService(ctrl),
		Contract:     mockResourcesv6.NewMockContractsService(ctrl),
		User:         mockResourcesv6.NewMockUsersService(ctrl),
		Group:        mockResourcesv6.NewMockGroupsService(ctrl),
		S3Key:        mockResourcesv6.NewMockS3KeysService(ctrl),
		BackupUnit:   mockResourcesv6.NewMockBackupUnitsService(ctrl),
		Pcc:          mockResourcesv6.NewMockPccsService(ctrl),
		K8s:          mockResourcesv6.NewMockK8sService(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocks) *CommandConfig {
	c.Locations = func() v6.LocationsService { return tm.Location }
	c.DataCenters = func() v6.DatacentersService { return tm.Datacenter }
	c.Servers = func() v6.ServersService { return tm.Server }
	c.Volumes = func() v6.VolumesService { return tm.Volume }
	c.Lans = func() v6.LansService { return tm.Lan }
	c.Nics = func() v6.NicsService { return tm.Nic }
	c.Loadbalancers = func() v6.LoadbalancersService { return tm.Loadbalancer }
	c.IpBlocks = func() v6.IpBlocksService { return tm.IpBlocks }
	c.Requests = func() v6.RequestsService { return tm.Request }
	c.Images = func() v6.ImagesService { return tm.Image }
	c.Snapshots = func() v6.SnapshotsService { return tm.Snapshot }
	c.FirewallRules = func() v6.FirewallRulesService { return tm.FirewallRule }
	c.Labels = func() v6.LabelResourcesService { return tm.Label }
	c.Contracts = func() v6.ContractsService { return tm.Contract }
	c.Users = func() v6.UsersService { return tm.User }
	c.Groups = func() v6.GroupsService { return tm.Group }
	c.S3Keys = func() v6.S3KeysService { return tm.S3Key }
	c.BackupUnit = func() v6.BackupUnitsService { return tm.BackupUnit }
	c.Pccs = func() v6.PccsService { return tm.Pcc }
	c.K8s = func() v6.K8sService { return tm.K8s }
	return c
}
