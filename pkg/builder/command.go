package builder

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Command struct {
	Command *cobra.Command

	childCommands []*Command
}

func (c *Command) AddCommand(commands ...*Command) {
	c.childCommands = append(c.childCommands, commands...)
	for _, cmd := range commands {
		c.Command.AddCommand(cmd.Command)
	}
}

func (c *Command) ChildCommands() []*Command {
	return c.childCommands
}

func (c *Command) AddStringFlag(name, shorthand, defaultValue, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringP(name, shorthand, defaultValue, desc)
	} else {
		flags.String(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.Command.Parent().Name(), c.Command.Name(), name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName("", c.Command.Name(), name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddStringSliceFlag(name, shorthand string, defaultValue []string, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringSliceP(name, shorthand, defaultValue, desc)
	} else {
		flags.StringSlice(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.Command.Parent().Name(), c.Command.Name(), name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName("", c.Command.Name(), name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddIntFlag(name, shorthand string, defaultValue int, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.IntP(name, shorthand, defaultValue, desc)
	} else {
		flags.Int(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.Command.Parent().Name(), c.Command.Name(), name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName("", c.Command.Name(), name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddFloat32Flag(name, shorthand string, defaultValue float32, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.Float32P(name, shorthand, defaultValue, desc)
	} else {
		flags.Float32(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.Command.Parent().Name(), c.Command.Name(), name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName("", c.Command.Name(), name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddBoolFlag(name, shorthand string, defaultValue bool, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.BoolP(name, shorthand, defaultValue, desc)
	} else {
		flags.Bool(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.Command.Parent().Name(), c.Command.Name(), name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName("", c.Command.Name(), name), c.Command.Flags().Lookup(name))
	}
}

func NewCommand(ctx context.Context, parent *Command, precr PreCommandRunner, cr CommandRunner, clitext, shortdesc, longdesc, example string, initServices bool) *Command {
	cc := &cobra.Command{
		Use:     clitext,
		Short:   shortdesc,
		Long:    longdesc,
		Example: example,
		PreRun: func(cmd *cobra.Command, args []string) {
			p := getPrinter()
			preCmdConfig := NewPreCommandConfig(
				p,
				clitext,
				getParentName(parent),
			)
			err := precr(preCmdConfig)
			clierror.CheckError(err, p.GetStderr())
		},
		Run: func(cmd *cobra.Command, args []string) {
			p := getPrinter()
			// Set Buffers
			cmd.SetIn(os.Stdin)
			cmd.SetOut(p.GetStdout())
			cmd.SetErr(p.GetStderr())

			cmdConfig, err := NewCommandConfig(
				ctx,
				os.Stdin,
				p,
				clitext,
				getParentName(parent),
				initServices,
			)
			clierror.CheckError(err, p.GetStderr())
			err = cr(cmdConfig)

			clierror.CheckError(err, p.GetStderr())
		},
	}
	c := &Command{
		Command: cc,
	}
	if parent != nil {
		parent.AddCommand(c)
	}

	return c
}

type PreCommandRunner func(*PreCommandConfig) error

type PreCommandConfig struct {
	Name       string
	ParentName string

	Printer printer.PrintService
}

func NewPreCommandConfig(p printer.PrintService, name, parentName string) *PreCommandConfig {
	return &PreCommandConfig{
		Name:       name,
		ParentName: parentName,
		Printer:    p,
	}
}

type CommandRunner func(*CommandConfig) error

type CommandConfig struct {
	Name       string
	ParentName string
	Stdin      io.Reader
	Printer    printer.PrintService
	Context    context.Context

	initServices func(*CommandConfig) error
	// Locations
	Locations func() resources.LocationsService
	// Services
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
}

func NewCommandConfig(ctx context.Context, in io.Reader, p printer.PrintService, name, parentName string, initServices bool) (*CommandConfig, error) {
	cmdConfig := &CommandConfig{
		Name:       name,
		ParentName: parentName,
		Stdin:      in,
		Printer:    p,
		Context:    ctx,

		initServices: func(c *CommandConfig) error {
			err := config.Load()
			if err != nil {
				return err
			}
			clientSvc, err := resources.NewClientService(
				viper.GetString(config.Username),
				viper.GetString(config.Password),
				viper.GetString(config.Token),
				viper.GetString(config.ArgServerUrl),
			)
			if err != nil {
				return err
			}

			c.Locations = func() resources.LocationsService { return resources.NewLocationService(clientSvc.Get(), c.Context) }
			c.DataCenters = func() resources.DatacentersService { return resources.NewDataCenterService(clientSvc.Get(), c.Context) }
			c.Servers = func() resources.ServersService { return resources.NewServerService(clientSvc.Get(), c.Context) }
			c.Volumes = func() resources.VolumesService { return resources.NewVolumeService(clientSvc.Get(), c.Context) }
			c.Lans = func() resources.LansService { return resources.NewLanService(clientSvc.Get(), c.Context) }
			c.Nics = func() resources.NicsService { return resources.NewNicService(clientSvc.Get(), c.Context) }
			c.Loadbalancers = func() resources.LoadbalancersService {
				return resources.NewLoadbalancerService(clientSvc.Get(), c.Context)
			}
			c.IpBlocks = func() resources.IpBlocksService { return resources.NewIpBlockService(clientSvc.Get(), c.Context) }
			c.Requests = func() resources.RequestsService { return resources.NewRequestService(clientSvc.Get(), c.Context) }
			c.Images = func() resources.ImagesService { return resources.NewImageService(clientSvc.Get(), c.Context) }
			c.Snapshots = func() resources.SnapshotsService { return resources.NewSnapshotService(clientSvc.Get(), c.Context) }
			c.FirewallRules = func() resources.FirewallRulesService {
				return resources.NewFirewallRuleService(clientSvc.Get(), c.Context)
			}
			c.Labels = func() resources.LabelResourcesService {
				return resources.NewLabelResourceService(clientSvc.Get(), c.Context)
			}
			c.Contracts = func() resources.ContractsService { return resources.NewContractService(clientSvc.Get(), c.Context) }
			return nil
		},
	}

	if initServices {
		if err := cmdConfig.initServices(cmdConfig); err != nil {
			return nil, err
		}
	}

	return cmdConfig, nil
}

func CheckRequiredGlobalFlags(parentCmdName string, globalFlagsName ...string) error {
	for _, flagName := range globalFlagsName {
		if viper.GetString(GetGlobalFlagName(parentCmdName, flagName)) == "" {
			return clierror.NewRequiredFlagErr(flagName)
		}
	}
	return nil
}

func CheckRequiredFlags(parentCmdName, cmdName string, localFlagsName ...string) error {
	for _, flagName := range localFlagsName {
		if viper.GetString(GetFlagName(parentCmdName, cmdName, flagName)) == "" {
			return clierror.NewRequiredFlagErr(flagName)
		}
	}
	return nil
}

func GetFlagName(prtName, cmdName, flagName string) string {
	if prtName != "" {
		return fmt.Sprintf("%s.%s.%s", prtName, cmdName, flagName)
	} else {
		return fmt.Sprintf("%s.%s", cmdName, flagName)
	}
}

func GetGlobalFlagName(cmdName, flagName string) string {
	return fmt.Sprintf("%s.%s", cmdName, flagName)
}

func getParentName(parent *Command) string {
	if parent != nil {
		return parent.Command.Name()
	} else {
		return ""
	}
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
