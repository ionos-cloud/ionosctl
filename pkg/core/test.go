package core

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	mockResourcesV5 "github.com/ionos-cloud/ionosctl/pkg/resources/v5/mocks"
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
	Client       *mockResourcesV5.MockClientService
	Location     *mockResourcesV5.MockLocationsService
	Datacenter   *mockResourcesV5.MockDatacentersService
	Server       *mockResourcesV5.MockServersService
	Volume       *mockResourcesV5.MockVolumesService
	Lan          *mockResourcesV5.MockLansService
	Nic          *mockResourcesV5.MockNicsService
	Loadbalancer *mockResourcesV5.MockLoadbalancersService
	IpBlocks     *mockResourcesV5.MockIpBlocksService
	Request      *mockResourcesV5.MockRequestsService
	Image        *mockResourcesV5.MockImagesService
	Snapshot     *mockResourcesV5.MockSnapshotsService
	FirewallRule *mockResourcesV5.MockFirewallRulesService
	Label        *mockResourcesV5.MockLabelResourcesService
	Contract     *mockResourcesV5.MockContractsService
	User         *mockResourcesV5.MockUsersService
	Group        *mockResourcesV5.MockGroupsService
	S3Key        *mockResourcesV5.MockS3KeysService
	BackupUnit   *mockResourcesV5.MockBackupUnitsService
	Pcc          *mockResourcesV5.MockPccsService
	K8s          *mockResourcesV5.MockK8sService
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
		Client:       mockResourcesV5.NewMockClientService(ctrl),
		Location:     mockResourcesV5.NewMockLocationsService(ctrl),
		Datacenter:   mockResourcesV5.NewMockDatacentersService(ctrl),
		Server:       mockResourcesV5.NewMockServersService(ctrl),
		Lan:          mockResourcesV5.NewMockLansService(ctrl),
		Volume:       mockResourcesV5.NewMockVolumesService(ctrl),
		Nic:          mockResourcesV5.NewMockNicsService(ctrl),
		Loadbalancer: mockResourcesV5.NewMockLoadbalancersService(ctrl),
		IpBlocks:     mockResourcesV5.NewMockIpBlocksService(ctrl),
		Request:      mockResourcesV5.NewMockRequestsService(ctrl),
		Image:        mockResourcesV5.NewMockImagesService(ctrl),
		Snapshot:     mockResourcesV5.NewMockSnapshotsService(ctrl),
		FirewallRule: mockResourcesV5.NewMockFirewallRulesService(ctrl),
		Label:        mockResourcesV5.NewMockLabelResourcesService(ctrl),
		Contract:     mockResourcesV5.NewMockContractsService(ctrl),
		User:         mockResourcesV5.NewMockUsersService(ctrl),
		Group:        mockResourcesV5.NewMockGroupsService(ctrl),
		S3Key:        mockResourcesV5.NewMockS3KeysService(ctrl),
		BackupUnit:   mockResourcesV5.NewMockBackupUnitsService(ctrl),
		Pcc:          mockResourcesV5.NewMockPccsService(ctrl),
		K8s:          mockResourcesV5.NewMockK8sService(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocks) *CommandConfig {
	c.Locations = func() v5.LocationsService { return tm.Location }
	c.DataCenters = func() v5.DatacentersService { return tm.Datacenter }
	c.Servers = func() v5.ServersService { return tm.Server }
	c.Volumes = func() v5.VolumesService { return tm.Volume }
	c.Lans = func() v5.LansService { return tm.Lan }
	c.Nics = func() v5.NicsService { return tm.Nic }
	c.Loadbalancers = func() v5.LoadbalancersService { return tm.Loadbalancer }
	c.IpBlocks = func() v5.IpBlocksService { return tm.IpBlocks }
	c.Requests = func() v5.RequestsService { return tm.Request }
	c.Images = func() v5.ImagesService { return tm.Image }
	c.Snapshots = func() v5.SnapshotsService { return tm.Snapshot }
	c.FirewallRules = func() v5.FirewallRulesService { return tm.FirewallRule }
	c.Labels = func() v5.LabelResourcesService { return tm.Label }
	c.Contracts = func() v5.ContractsService { return tm.Contract }
	c.Users = func() v5.UsersService { return tm.User }
	c.Groups = func() v5.GroupsService { return tm.Group }
	c.S3Keys = func() v5.S3KeysService { return tm.S3Key }
	c.BackupUnit = func() v5.BackupUnitsService { return tm.BackupUnit }
	c.Pccs = func() v5.PccsService { return tm.Pcc }
	c.K8s = func() v5.K8sService { return tm.K8s }
	return c
}
