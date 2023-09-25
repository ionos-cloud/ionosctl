package core

import (
	"context"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/v6/services/certmanager"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	container_registry "github.com/ionos-cloud/ionosctl/v6/services/container-registry"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/spf13/cobra"
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
			cmd.SilenceUsage = true

			// Set Command to Command Builder
			// The cmd is passed to the PreCommandCfg
			info.Command = &Command{Command: cmd}
			// Create New PreCommandCfg
			preCmdConfig := NewPreCommandCfg(info)
			err := info.PreCmdRun(preCmdConfig)
			if err != nil {
				die.Die(err.Error())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// Set Command to Command Builder
			// The cmd is passed to the CommandCfg
			info.Command = &Command{Command: cmd}
			// Create New CommandCfg
			cmdConfig, err := NewCommandCfg(ctx, info)
			if err != nil {
				die.Die(err.Error())
			}

			err = info.CmdRun(cmdConfig)
			if err != nil {
				die.Die(err.Error())
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
}

func NewPreCommandCfg(info CommandBuilder) *PreCommandConfig {
	return &PreCommandConfig{
		Command:   info.Command,
		NS:        info.GetNS(),
		Namespace: info.Namespace,
		Resource:  info.Resource,
		Verb:      info.Verb,
	}
}

func NewCommandCfg(ctx context.Context, info CommandBuilder) (*CommandConfig, error) {
	cmdConfig := &CommandConfig{
		Command:   info.Command,
		NS:        info.GetNS(),
		Namespace: info.Namespace,
		Resource:  info.Resource,
		Verb:      info.Verb,
		Stdin:     info.Command.Command.InOrStdin(),
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
