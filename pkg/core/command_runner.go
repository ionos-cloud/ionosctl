package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/v6/services/certmanager"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	container_registry "github.com/ionos-cloud/ionosctl/v6/services/container-registry"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand(ctx context.Context, parent *Command, info CommandBuilder) *Command {
	if info.PreCmdRun == nil {
		// Fixes a nil pointer reference panic
		info.PreCmdRun = NoPreRun
	}

	cc := &cobra.Command{
		Use:     info.Verb,
		Short:   info.ShortDesc,
		Long:    info.LongDesc,
		Aliases: info.Aliases,
		Example: info.Example,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Set Printer in sync with the Output Flag
			noHeaders, _ := cmd.Flags().GetBool(constants.ArgNoHeaders)
			p := getPrinter(noHeaders)
			// Set Command to Command Builder
			// The cmd is passed to the PreCommandCfg
			info.Command = &Command{Command: cmd}
			// Create New PreCommandCfg
			preCmdConfig := NewPreCommandCfg(p, info)
			err := info.PreCmdRun(preCmdConfig)
			if err != nil {
				return fmt.Errorf("prerun error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set Printer in sync with the Output Flag
			noHeaders, _ := cmd.Flags().GetBool(constants.ArgNoHeaders)
			p := getPrinter(noHeaders)
			// Set Buffers
			cmd.SetIn(os.Stdin)
			cmd.SetOut(p.GetStdout())
			cmd.SetErr(p.GetStderr())
			// Set Command to Command Builder
			// The cmd is passed to the CommandCfg
			info.Command = &Command{Command: cmd}
			// Create New CommandCfg
			cmdConfig, err := NewCommandCfg(ctx, os.Stdin, p, info)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
			err = info.CmdRun(cmdConfig)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
			return nil
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

var NoPreRun = func(c *PreCommandConfig) error { return nil }

// PreCommandConfig Properties
type PreCommandConfig struct {
	// Command is a Wrapper around Cobra Command
	Command *Command

	// NS is Global Namespace for all Command Levels
	NS string
	// Namespace is the first level of the Command. e.g. [ionosctl] server
	Namespace string
	// Resource is the second level of the Command. e.g. [ionosctl server] volume
	Resource string
	// Verb is the 3rd level of the Command. e.g. [ionosctl server volume] attach
	Verb string

	// Printer used in output formatting
	Printer printer.PrintService
}

func NewPreCommandCfg(p printer.PrintService, info CommandBuilder) *PreCommandConfig {
	return &PreCommandConfig{
		Command:   info.Command,
		NS:        info.GetNS(),
		Namespace: info.Namespace,
		Resource:  info.Resource,
		Verb:      info.Verb,
		Printer:   p,
	}
}

func NewCommandCfg(ctx context.Context, in io.Reader, p printer.PrintService, info CommandBuilder) (
	*CommandConfig, error,
) {
	cmdConfig := &CommandConfig{
		Command:   info.Command,
		NS:        info.GetNS(),
		Namespace: info.Namespace,
		Resource:  info.Resource,
		Verb:      info.Verb,
		Stdin:     in,
		Printer:   p,
		Context:   ctx,
		// Define cmd Command Config function for Command
		initCfg: func(c *CommandConfig) error {
			// Init Clients and Services
			client, err := client2.Get()
			if err != nil {
				return err
			}
			if err := c.CloudApiV6Services.InitServices(client); err != nil {
				return err
			}
			if err = c.AuthV1Services.InitServices(client); err != nil {
				return err
			}

			if err = c.CloudApiDbaasPgsqlServices.InitServices(client); err != nil {
				return err
			}

			if err = c.CertificateManagerServices.InitServices(client); err != nil {
				return err
			}

			if err = c.ContainerRegistryServices.InitServices(client); err != nil {
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
	// Command is a Wrapper around Cobra Command
	Command *Command

	// CLI Engine
	// NS is Global Namespace for all Command Levels
	NS string
	// Namespace is the first level of the Command. e.g. [ionosctl] server
	Namespace string
	// Resource is the second level of the Command. e.g. [ionosctl server] volume
	Resource string
	// Verb is the 3rd level of the Command. e.g. [ionosctl server volume] attach
	Verb    string
	Stdin   io.Reader
	Printer printer.PrintService
	initCfg func(commandConfig *CommandConfig) error

	// Services
	CloudApiV6Services         cloudapiv6.Services
	AuthV1Services             authservice.Services
	CloudApiDbaasPgsqlServices cloudapidbaaspgsql.Services
	CertificateManagerServices certmanager.Services
	ContainerRegistryServices  container_registry.Services

	// Context
	Context context.Context
}

// TODO: Seems like there's no better way to Verbose print outside of 'commands' pkg, other than instantiating a PrintService as so. PrintService merits a refactor. It seems like without this exported func, I can only make Verbose prints if I am inside of a `commands` command object.
func GetPrinter(noHeaders bool) printer.PrintService {
	return getPrinter(noHeaders)
}

func getPrinter(noHeaders bool) printer.PrintService {
	var out io.Writer
	if viper.GetBool(constants.ArgQuiet) {
		var execOut bytes.Buffer
		out = &execOut
	} else {
		out = os.Stdout // lol we should either not allow CommandBuilder to customize out buffer at all, or find a way for it to influence this line. I can't change command output in tests because of this
	}
	printReg, err := printer.NewPrinterRegistry(out, os.Stderr, noHeaders)
	clierror.CheckErrorAndDie(err, os.Stderr)
	return printReg[viper.GetString(constants.ArgOutput)]
}
