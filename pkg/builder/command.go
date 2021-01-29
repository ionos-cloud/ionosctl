package builder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
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
	viper.BindPFlag(GetFlagName(c.Command.Name(), name), c.Command.Flags().Lookup(name))
}

func (c *Command) AddBoolFlag(name, shorthand string, defaultValue bool, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.BoolP(name, shorthand, defaultValue, desc)
	} else {
		flags.Bool(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.Command.Name(), name), c.Command.Flags().Lookup(name))
}

func NewCommand(ctx context.Context, parent *Command, cr CommandRunner, clitext, shortdesc, longdesc string, initServices bool) *Command {
	cc := &cobra.Command{
		Use:   clitext,
		Short: shortdesc,
		Long:  longdesc,
		Run: func(cmd *cobra.Command, args []string) {
			var out io.Writer
			if viper.GetBool(config.ArgQuiet) {
				var execOut bytes.Buffer
				out = &execOut
			} else {
				out = os.Stdout
			}
			printReg := utils.NewPrinterRegistry(out, os.Stderr)
			printer := printReg[viper.GetString(config.ArgOutput)]
			// Set Buffers
			cmd.SetIn(os.Stdin)
			cmd.SetOut(printer.GetStdout())
			cmd.SetErr(printer.GetStderr())

			cmdConfig, err := NewCommandConfig(
				ctx,
				os.Stdin,
				printer,
				clitext,
				initServices,
			)
			utils.CheckError(err, printer.GetStderr())
			err = cr(cmdConfig)

			utils.CheckError(err, printer.GetStderr())
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

type CommandRunner func(*CommandConfig) error

type CommandConfig struct {
	Name  string
	Stdin io.Reader

	Printer utils.PrintService
	Context context.Context

	initServices func(*CommandConfig) error
	// Services
	DataCenters func() resources.DatacentersService
}

func NewCommandConfig(ctx context.Context, in io.Reader, p utils.PrintService, name string, initServices bool, ) (*CommandConfig, error) {
	cmdConfig := &CommandConfig{
		Name:    name,
		Stdin:   in,
		Printer: p,

		initServices: func(c *CommandConfig) error {
			err := config.LoadFile()
			if err != nil {
				return err
			}
			clientSvc, err := resources.NewClientService(
				viper.GetString(config.Username),
				viper.GetString(config.Password),
				viper.GetString(config.ArgServerUrl),
			)
			if err != nil {
				return err
			}

			c.DataCenters = func() resources.DatacentersService { return resources.NewDataCenterService(clientSvc.Get(), ctx) }
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

func GetFlagName(cmdName, name string) string {
	return fmt.Sprintf("%s.%s", cmdName, name)
}
