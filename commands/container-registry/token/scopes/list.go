package scopes

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenScopesListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "token",
			Resource:   "scope",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List a token's scopes",
			LongDesc:   "List all scopes of a token, showing for each its ScopeId (zero-based index used by 'scope delete'), name (target resource), type and allowed actions.",
			Example:    "ionosctl container-registry token scope list --registry-id REGISTRY_ID --token-id TOKEN_ID",
			PreCmdRun:  PreCmdTokenScopesList,
			CmdRun:     CmdGetTokenScopesList,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry that owns the token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "", "", "The unique ID of the token whose scopes to list")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId,
		func(cobracmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(
				viper.GetString(
					core.GetFlagName(
						cmd.NS, constants.FlagRegistryId,
					),
				),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdGetTokenScopesList(c *core.CommandConfig) error {
	reg_id := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	token_id := viper.GetString(core.GetFlagName(c.NS, FlagTokenId))

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensFindById(context.Background(), reg_id, token_id).Execute()
	if err != nil {
		return err
	}

	properties, ok := token.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("could not retrieve Container Registry Token properties")
	}

	scopes, ok := properties.GetScopesOk()
	if !ok || scopes == nil {
		return fmt.Errorf("could not retrieve Container Registry Token Scopes")
	}

	t := table.New(allScopeCols)
	if err := t.Extract(token.Properties.Scopes); err != nil {
		return err
	}
	for i := range token.Properties.Scopes {
		t.SetCell(i, "ScopeId", fmt.Sprintf("%v", i))
	}

	return c.Out(t.Render(table.ResolveCols(allScopeCols, c.Cols())))
}

func PreCmdTokenScopesList(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagTokenId, constants.FlagRegistryId)
	if err != nil {
		return err
	}
	return nil
}
