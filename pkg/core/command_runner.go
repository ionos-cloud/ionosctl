package core

import (
	"bytes"
	"context"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
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
			client, err := c.InitClient()
			if err != nil {
				return err
			}
			if err = c.InitServices(client); err != nil {
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
	// Resources Services
	BackupUnit    func() resources.BackupUnitsService
	Contracts     func() resources.ContractsService
	DataCenters   func() resources.DatacentersService
	FirewallRules func() resources.FirewallRulesService
	Groups        func() resources.GroupsService
	Images        func() resources.ImagesService
	IpBlocks      func() resources.IpBlocksService
	K8s           func() resources.K8sService
	Labels        func() resources.LabelResourcesService
	Lans          func() resources.LansService
	Loadbalancers func() resources.LoadbalancersService
	Locations     func() resources.LocationsService
	Nics          func() resources.NicsService
	Pccs          func() resources.PccsService
	Requests      func() resources.RequestsService
	S3Keys        func() resources.S3KeysService
	Servers       func() resources.ServersService
	Snapshots     func() resources.SnapshotsService
	Templates     func() resources.TemplatesService
	Users         func() resources.UsersService
	Volumes       func() resources.VolumesService
	// Context
	Context context.Context
}

// InitClient for Commands
func (c *CommandConfig) InitClient() (*resources.Client, error) {
	err := config.Load()
	if err != nil {
		return nil, err
	}
	clientSvc, err := resources.NewClientService(
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

// InitServices for Commands
func (c *CommandConfig) InitServices(client *resources.Client) error {
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
	c.Templates = func() resources.TemplatesService { return resources.NewTemplateService(client, c.Context) }
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
