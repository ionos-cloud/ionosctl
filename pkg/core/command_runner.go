package core

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
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
			v5Client, err := c.InitV5Client()
			if err != nil {
				return err
			}
			if err = c.InitV5Services(v5Client); err != nil {
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
	// V5 Resources Services
	Locations     func() v5.LocationsService
	DataCenters   func() v5.DatacentersService
	Servers       func() v5.ServersService
	Volumes       func() v5.VolumesService
	Lans          func() v5.LansService
	Nics          func() v5.NicsService
	Loadbalancers func() v5.LoadbalancersService
	Requests      func() v5.RequestsService
	Images        func() v5.ImagesService
	Snapshots     func() v5.SnapshotsService
	IpBlocks      func() v5.IpBlocksService
	FirewallRules func() v5.FirewallRulesService
	Labels        func() v5.LabelResourcesService
	Contracts     func() v5.ContractsService
	Users         func() v5.UsersService
	Groups        func() v5.GroupsService
	S3Keys        func() v5.S3KeysService
	BackupUnit    func() v5.BackupUnitsService
	Pccs          func() v5.PccsService
	K8s           func() v5.K8sService
	// Context
	Context context.Context
}

// InitV5Client for Commands
func (c *CommandConfig) InitV5Client() (*v5.Client, error) {
	err := config.Load()
	if err != nil {
		return nil, err
	}
	clientSvc, err := v5.NewClientService(
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

// InitV5Services for Commands
func (c *CommandConfig) InitV5Services(client *v5.Client) error {
	c.Locations = func() v5.LocationsService { return v5.NewLocationService(client, c.Context) }
	c.DataCenters = func() v5.DatacentersService { return v5.NewDataCenterService(client, c.Context) }
	c.Servers = func() v5.ServersService { return v5.NewServerService(client, c.Context) }
	c.Volumes = func() v5.VolumesService { return v5.NewVolumeService(client, c.Context) }
	c.Lans = func() v5.LansService { return v5.NewLanService(client, c.Context) }
	c.Nics = func() v5.NicsService { return v5.NewNicService(client, c.Context) }
	c.Loadbalancers = func() v5.LoadbalancersService { return v5.NewLoadbalancerService(client, c.Context) }
	c.IpBlocks = func() v5.IpBlocksService { return v5.NewIpBlockService(client, c.Context) }
	c.Requests = func() v5.RequestsService { return v5.NewRequestService(client, c.Context) }
	c.Images = func() v5.ImagesService { return v5.NewImageService(client, c.Context) }
	c.Snapshots = func() v5.SnapshotsService { return v5.NewSnapshotService(client, c.Context) }
	c.FirewallRules = func() v5.FirewallRulesService { return v5.NewFirewallRuleService(client, c.Context) }
	c.Labels = func() v5.LabelResourcesService { return v5.NewLabelResourceService(client, c.Context) }
	c.Contracts = func() v5.ContractsService { return v5.NewContractService(client, c.Context) }
	c.Users = func() v5.UsersService { return v5.NewUserService(client, c.Context) }
	c.Groups = func() v5.GroupsService { return v5.NewGroupService(client, c.Context) }
	c.S3Keys = func() v5.S3KeysService { return v5.NewS3KeyService(client, c.Context) }
	c.BackupUnit = func() v5.BackupUnitsService { return v5.NewBackupUnitService(client, c.Context) }
	c.Pccs = func() v5.PccsService { return v5.NewPrivateCrossConnectService(client, c.Context) }
	c.K8s = func() v5.K8sService { return v5.NewK8sService(client, c.Context) }
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
