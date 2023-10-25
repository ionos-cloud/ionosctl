package core

import (
	"context"
	"fmt"
	"io"
	"os"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/v6/services/certmanager"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	container_registry "github.com/ionos-cloud/ionosctl/v6/services/container-registry"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommandWithJsonProperties(ctx context.Context, parent *Command, jsonExample string, toUnmarshal interface{}, info CommandBuilder) *Command {
	newInfo := info
	// Inject custom behaviour for the command run
	newInfo.CmdRun = withJsonFile(jsonExample, toUnmarshal, info.CmdRun)
	// Inject custom flags
	cmd := NewCommand(ctx, parent, newInfo)

	cmd.Command.Flags().String(constants.FlagJsonProperties, "",
		"Path to a JSON file containing the desired properties. Overrides any other properties set.")
	viper.BindPFlag(constants.FlagJsonProperties, cmd.Command.Flags().Lookup(constants.FlagJsonProperties))

	if jsonExample != "" {
		cmd.Command.Flags().Bool(constants.FlagJsonPropertiesExample, false,
			fmt.Sprintf("If set, prints a complete JSON which could be used for --%s "+
				"and exits. Hint: Pipe me to a .json file", constants.FlagJsonProperties))
		viper.BindPFlag(constants.FlagJsonPropertiesExample, cmd.Command.Flags().Lookup(constants.FlagJsonPropertiesExample))
	}

	viper.BindPFlag(constants.FlagJsonProperties, cmd.Command.Flags().Lookup(constants.FlagJsonProperties))
	viper.BindPFlag(constants.FlagJsonPropertiesExample, cmd.Command.Flags().Lookup(constants.FlagJsonPropertiesExample))

	return cmd
}

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
				return err
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
				return err
			}

			err = info.CmdRun(cmdConfig)
			if err != nil {
				return err
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

// withJsonFile decorates your 'run' CommandRun to add functionality related to the --json-properties flag.
// If the flag is provided, the json file is read and unmarshalled into the toUnmarshal interface
// you should give a pointer to an SDK struct as a parameter to `toUnmarshal`
//
// 'example' is a string containing a JSON example of the struct you should pass to `toUnmarshal`,
// which the user can print with --json-properties-example (and be able to pipe into a .json file, for instance)
// if 'example' is empty, this functionality is omitted
//
// `run` is the CommandRun you want to decorate - if the flag is not provided, `run` will be called as usual
func withJsonFile(example string, toUnmarshal interface{}, run CommandRun) CommandRun {
	return func(c *CommandConfig) error {
		printExample, err := c.Command.Command.Flags().GetBool(constants.FlagJsonPropertiesExample)
		if err != nil {
			return err
		}
		if printExample {
			fmt.Fprintf(c.Command.Command.OutOrStdout(), example)
			return nil
		}
		jsonFile, err := c.Command.Command.Flags().GetString(constants.FlagJsonProperties)
		if err != nil {
			return err
		}

		// Check if the json-properties flag is provided
		if jsonFile == "" {
			// No JSON properties file specified, directly run the command.
			// The inner command can handle manual flag parsing.
			return run(c)
		}

		// Check if the file actually exists
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			return fmt.Errorf("specified json properties file does not exist: %s", jsonFile)
		}

		v := viper.New()
		v.SetConfigFile(jsonFile)
		err = v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", jsonFile, err)
		}

		// -- If properties key exists and is a map, promote its key-value pairs
		// Unmarshal the config into a map
		// var configMap map[string]interface{}
		// err = v.Unmarshal(&configMap)
		// if err != nil {
		// 	return fmt.Errorf("failed unmarshalling config: %w", err)
		// }

		// if propValue, exists := configMap["properties"].(map[string]interface{}); exists {
		// 	delete(configMap, "properties") // Delete the properties key
		// 	for k, val := range propValue {
		// 		v.Set(k, val)
		// 	}
		// }

		err = v.Unmarshal(toUnmarshal)
		if err != nil {
			return fmt.Errorf("failed unmarshalling json properties into object: %w", err)
		}

		return run(c)
	}
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

func withJsonFile(example string, toUnmarshal interface{}, run CommandRun) CommandRun {
	return func(c *CommandConfig) error {
		if viper.GetBool(constants.FlagJsonPropertiesExample) {
			fmt.Fprintf(c.Command.Command.OutOrStdout(), example)
			return nil
		}
		jsonFile := viper.GetString(constants.FlagJsonProperties)

		// Check if the json-properties flag is provided
		if jsonFile == "" {
			// No JSON properties file specified, directly run the command.
			// The inner command can handle manual flag parsing.
			return run(c)
		}

		// Check if the file actually exists
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			return fmt.Errorf("specified json properties file does not exist: %s", jsonFile)
		}

		v := viper.New()
		v.SetConfigFile(jsonFile)
		err := v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed reading %s: %w", jsonFile, err)
		}

		// Unmarshal the config into a map
		var configMap map[string]interface{}
		err = v.Unmarshal(&configMap)
		if err != nil {
			return fmt.Errorf("failed unmarshalling config: %w", err)
		}

		// If properties key exists and is a map, promote its key-value pairs
		if propValue, exists := configMap["properties"].(map[string]interface{}); exists {
			delete(configMap, "properties") // Delete the properties key
			for k, val := range propValue {
				v.Set(k, val)
			}
		}

		err = v.Unmarshal(toUnmarshal)
		if err != nil {
			return fmt.Errorf("failed unmarshalling json properties into object: %w", err)
		}

		return run(c)
	}
}

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
