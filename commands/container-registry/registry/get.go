package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func RegGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "get",
			Aliases:    []string{"g"},
			ShortDesc:  "Get the properties of a registry",
			LongDesc:   "Get the properties of a single registry, including its hostname (used for 'docker login'), location, garbage-collection schedule, enabled features, and provisioning state.",
			Example:    "ionosctl container-registry registry get --id REGISTRY_ID",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, "i", "", "The unique ID of the registry to retrieve", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdGet(c *core.CommandConfig) error {
	id, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	reg, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesFindById(context.Background(), id).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(reg)
}

func PreCmdGet(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
