package core

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(ctx context.Context, parent *Command, info CommandBuilder) *Command {
	cc := &cobra.Command{
		Use:     info.Verb,
		Short:   info.ShortDesc,
		Long:    info.LongDesc,
		Aliases: info.Aliases,
		Example: info.Example,
		PreRun: func(cmd *cobra.Command, args []string) {
			p := getPrinter()
			preCmdConfig := NewPreCommandCfg(p, info)
			err := info.PreCmdRun(preCmdConfig)
			clierror.CheckError(err, p.GetStderr())
		},
		Run: func(cmd *cobra.Command, args []string) {
			p := getPrinter()
			// Set Buffers
			cmd.SetIn(os.Stdin)
			cmd.SetOut(p.GetStdout())
			cmd.SetErr(p.GetStderr())
			cmdConfig, err := NewCommandCfg(ctx, os.Stdin, p, info)
			clierror.CheckError(err, p.GetStderr())
			err = info.CmdRun(cmdConfig)
			clierror.CheckError(err, p.GetStderr())
		},
	}
	c := &Command{
		NS:      info.GetNS(),
		Command: cc,
	}
	if parent != nil {
		parent.AddCommand(c)
	}

	return c
}

// PreCommandRun will run in PreRun of Cobra Command structure, before running the actual Command.
// Its purpose is to keep the validate part separate from run part.
type PreCommandRun func(commandConfig *PreCommandConfig) error

// PreCommandConfig Properties
type PreCommandConfig struct {
	// NS is Global Namespace for all Command Levels
	NS string
	// Namespace is the first level of the Command. e.g. [ionosctl] server
	Namespace string
	// Resource is the second level of the Command. e.g. [ionosctl server] volume
	Resource string
	// Verb is the 3rd level of the Command. e.g. [ionosctl server volume] attach
	Verb string

	Printer printer.PrintService
}

func NewPreCommandCfg(p printer.PrintService, info CommandBuilder) *PreCommandConfig {
	return &PreCommandConfig{
		NS:        info.GetNS(),
		Namespace: info.Namespace,
		Resource:  info.Resource,
		Verb:      info.Verb,
		Printer:   p,
	}
}

func NewCommandCfg(ctx context.Context, in io.Reader, p printer.PrintService, info CommandBuilder) (*CommandConfig, error) {
	cmdConfig := &CommandConfig{
		NS:        info.GetNS(),
		Namespace: info.Namespace,
		Resource:  info.Resource,
		Verb:      info.Verb,
		Stdin:     in,
		Printer:   p,
		Context:   ctx,
		// Define init Command Config function for Command
		initCfg: func(c *CommandConfig) error {
			client, err := c.InitV6Client()
			if err != nil {
				return err
			}
			if err = c.InitV6Services(client); err != nil {
				return err
			}
			return nil
		},
	}
	if info.InitClient {
		err := cmdConfig.initCfg(cmdConfig)
		if err != nil {
			return nil, err
		}
	}
	return cmdConfig, nil
}

type CommandRun func(commandConfig *CommandConfig) error

// CommandConfig Properties and Services
type CommandConfig struct {
	// NS is Global Namespace for all Command Levels
	NS string
	// Namespace is the first level of the Command. e.g. [ionosctl] server
	Namespace string
	// Resource is the second level of the Command. e.g. [ionosctl server] volume
	Resource string
	// Verb is the 3rd level of the Command. e.g. [ionosctl server volume] attach
	Verb string

	Stdin   io.Reader
	Printer printer.PrintService
	initCfg func(commandConfig *CommandConfig) error
	// V6 Resources Services
	Locations            func() v6.LocationsService
	DataCenters          func() v6.DatacentersService
	Servers              func() v6.ServersService
	Volumes              func() v6.VolumesService
	Lans                 func() v6.LansService
	NatGateways          func() v6.NatGatewaysService
	NetworkLoadBalancers func() v6.NetworkLoadBalancersService
	Nics                 func() v6.NicsService
	Loadbalancers        func() v6.LoadbalancersService
	Requests             func() v6.RequestsService
	Images               func() v6.ImagesService
	Snapshots            func() v6.SnapshotsService
	IpBlocks             func() v6.IpBlocksService
	FirewallRules        func() v6.FirewallRulesService
	FlowLogs             func() v6.FlowLogsService
	Labels               func() v6.LabelResourcesService
	Contracts            func() v6.ContractsService
	Users                func() v6.UsersService
	Groups               func() v6.GroupsService
	S3Keys               func() v6.S3KeysService
	BackupUnit           func() v6.BackupUnitsService
	Pccs                 func() v6.PccsService
	K8s                  func() v6.K8sService
	Templates            func() v6.TemplatesService
	// Context
	Context context.Context
}

// InitV6Client for Commands
func (c *CommandConfig) InitV6Client() (*v6.Client, error) {
	err := config.Load()
	if err != nil {
		return nil, err
	}
	clientSvc, err := v6.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token), // Token support
		viper.GetString(config.ArgServerUrl),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}

// InitV6Services for Commands
func (c *CommandConfig) InitV6Services(client *v6.Client) error {
	c.Locations = func() v6.LocationsService { return v6.NewLocationService(client, c.Context) }
	c.DataCenters = func() v6.DatacentersService { return v6.NewDataCenterService(client, c.Context) }
	c.Servers = func() v6.ServersService { return v6.NewServerService(client, c.Context) }
	c.Volumes = func() v6.VolumesService { return v6.NewVolumeService(client, c.Context) }
	c.Lans = func() v6.LansService { return v6.NewLanService(client, c.Context) }
	c.NatGateways = func() v6.NatGatewaysService { return v6.NewNatGatewayService(client, c.Context) }
	c.NetworkLoadBalancers = func() v6.NetworkLoadBalancersService {
		return v6.NewNetworkLoadBalancerService(client, c.Context)
	}
	c.Nics = func() v6.NicsService { return v6.NewNicService(client, c.Context) }
	c.Loadbalancers = func() v6.LoadbalancersService { return v6.NewLoadbalancerService(client, c.Context) }
	c.IpBlocks = func() v6.IpBlocksService { return v6.NewIpBlockService(client, c.Context) }
	c.Requests = func() v6.RequestsService { return v6.NewRequestService(client, c.Context) }
	c.Images = func() v6.ImagesService { return v6.NewImageService(client, c.Context) }
	c.Snapshots = func() v6.SnapshotsService { return v6.NewSnapshotService(client, c.Context) }
	c.FirewallRules = func() v6.FirewallRulesService { return v6.NewFirewallRuleService(client, c.Context) }
	c.FlowLogs = func() v6.FlowLogsService { return v6.NewFlowLogService(client, c.Context) }
	c.Labels = func() v6.LabelResourcesService { return v6.NewLabelResourceService(client, c.Context) }
	c.Contracts = func() v6.ContractsService { return v6.NewContractService(client, c.Context) }
	c.Users = func() v6.UsersService { return v6.NewUserService(client, c.Context) }
	c.Groups = func() v6.GroupsService { return v6.NewGroupService(client, c.Context) }
	c.S3Keys = func() v6.S3KeysService { return v6.NewS3KeyService(client, c.Context) }
	c.BackupUnit = func() v6.BackupUnitsService { return v6.NewBackupUnitService(client, c.Context) }
	c.Pccs = func() v6.PccsService { return v6.NewPrivateCrossConnectService(client, c.Context) }
	c.K8s = func() v6.K8sService { return v6.NewK8sService(client, c.Context) }
	c.Templates = func() v6.TemplatesService { return v6.NewTemplateService(client, c.Context) }
	return nil
}

func getPrinter() printer.PrintService {
	var out io.Writer
	if viper.GetBool(config.ArgQuiet) {
		var execOut bytes.Buffer
		out = &execOut
	} else {
		out = os.Stdout
	}
	printReg, err := printer.NewPrinterRegistry(out, os.Stderr)
	clierror.CheckError(err, os.Stderr)
	return printReg[viper.GetString(config.ArgOutput)]
}
