package scopes

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenScopesDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "token",
			Resource:  "scope",
			Verb:      "delete",
			Aliases:   []string{"d", "rm", "remove"},
			ShortDesc: "Delete a token scope",
			LongDesc: "Use this command to delete a token scope of a container registry. If a name is provided, the first scope with that" +
				" name will be deleted. It is possible to delete all scopes by providing the --all flag.",
			Example:    "ionosctl container-registry token scope delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]",
			PreCmdRun:  PreCmdTokenScopesDelete,
			CmdRun:     CmdGetTokenScopesDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "t", "", "Token ID")
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

	cmd.AddIntFlag(FlagScopeId, "n", -1, "Scope id")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all scopes of all tokens of a registry.")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allScopeCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allScopeCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdGetTokenScopesDelete(c *core.CommandConfig) error {
	regId, _ := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	tokenId, _ := c.Command.Command.Flags().GetString(FlagTokenId)

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensFindById(context.Background(), regId, tokenId).Execute()
	if err != nil {
		return err
	}

	allFlag, _ := c.Command.Command.Flags().GetBool(constants.ArgAll)
	if allFlag {
		updateToken := containerregistry.NewPutTokenInputWithDefaults()
		updateProp := containerregistry.NewPostTokenPropertiesWithDefaults()
		if token.Properties.ExpiryDate != nil {
			updateProp.SetExpiryDate(token.Properties.GetExpiryDate())
		}
		updateProp.SetStatus(token.Properties.GetStatus())
		updateProp.SetName(token.Properties.GetName())
		updateToken.SetProperties(*updateProp)

		msg := fmt.Sprintf("delete all scopes from Token: %s", *token.Id)

		if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = client.Must().RegistryClient.TokensApi.RegistriesTokensDelete(context.Background(), regId, tokenId).Execute()
		if err != nil {
			return err
		}

		_, _, err = client.Must().RegistryClient.TokensApi.RegistriesTokensPut(context.Background(), regId, tokenId).PutTokenInput(*updateToken).Execute()
		if err != nil {
			return err
		}

		return nil
	}

	id, err := c.Command.Command.Flags().GetInt(FlagScopeId)
	if err != nil {
		return err
	}

	scopes := token.Properties.Scopes
	if id < 0 || id >= len(scopes) {
		return fmt.Errorf("invalid scope ID %d: out of range", id)
	}

	updateToken := containerregistry.NewPutTokenInputWithDefaults()
	updateProp := containerregistry.NewPostTokenPropertiesWithDefaults()

	scopes = append(scopes[:id], scopes[id+1:]...)

	if token.Properties.ExpiryDate != nil {
		updateProp.SetExpiryDate(token.Properties.GetExpiryDate())
	}
	updateProp.SetStatus(token.Properties.GetStatus())
	updateProp.SetName(token.Properties.GetName())
	updateProp.SetScopes(scopes)
	updateToken.SetProperties(*updateProp)

	targetScope := token.Properties.Scopes[id]
	ask := fmt.Sprintf("delete scope %d (name '%s', type '%s' with actions [%s]) from Token: %s", id,
		targetScope.Name, targetScope.Type, strings.Join(targetScope.Actions, ", "), token.Properties.Name)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), ask, viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err = client.Must().RegistryClient.TokensApi.RegistriesTokensDelete(context.Background(), regId, tokenId).Execute()
	if err != nil {
		return err
	}

	_, _, err = client.Must().RegistryClient.TokensApi.RegistriesTokensPut(context.Background(), regId, tokenId).PutTokenInput(*updateToken).Execute()
	if err != nil {
		return err
	}

	return nil

}

func PreCmdTokenScopesDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagRegistryId, FlagTokenId, FlagScopeId},
		[]string{constants.FlagRegistryId, FlagTokenId, constants.ArgAll},
	)
}
