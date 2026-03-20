package token

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List all tokens",
			LongDesc:   "List all tokens for your container registry",
			Example:    "ionosctl container-registry token list --registry-id [REGISTRY-ID]",
			PreCmdRun:  PreCmdListToken,
			CmdRun:     CmdListToken,
			InitClient: true,
		},
	)

	cmd.AddBoolFlag(constants.ArgAll, "a", false, "List all tokens, including expired ones")
	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdListToken(c *core.CommandConfig) error {
	allFlag := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if !allFlag {
		id := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))

		tokens, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), id).Execute()
		if err != nil {
			return err
		}
		return c.Out(table.Sprint(allCols, tokens, cols, table.WithPrefix("items")))
	}

	var list = make([]containerregistry.TokensResponse, 0)

	regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	if items, ok := regs.GetItemsOk(); ok && items != nil {
		for _, reg := range items {
			tokens, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensGet(context.Background(), *reg.Id).Execute()
			if err != nil {
				return err
			}

			list = append(list, tokens)
		}
	}

	return c.Out(table.Sprint(allCols, list, cols, table.WithPrefix("*.items")))
}

func PreCmdListToken(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagRegistryId},
		[]string{constants.ArgAll},
	)
}
