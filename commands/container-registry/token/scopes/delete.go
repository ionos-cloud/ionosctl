package scopes

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
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

	cmd.AddStringFlag(FlagRegId, "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTokenId, "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddIntFlag(FlagScopeId, "n", -1, "Scope id")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all scopes of all tokens of a registry.")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allScopeCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allScopeCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdGetTokenScopesDelete(c *core.CommandConfig) error {
	regId := viper.GetString(core.GetFlagName(c.NS, FlagRegId))
	tokenId := viper.GetString(core.GetFlagName(c.NS, FlagTokenId))
	token, _, err := c.ContainerRegistryServices.Token().Get(tokenId, regId)
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		updateToken := sdkgo.NewPutTokenInputWithDefaults()
		updateProp := sdkgo.NewPostTokenPropertiesWithDefaults()

		if token.Properties.GetExpiryDate() != nil {
			updateProp.SetExpiryDate(*token.Properties.GetExpiryDate())
		}
		if token.Properties.GetStatus() != nil {
			updateProp.SetStatus(*token.Properties.GetStatus())
		}
		updateProp.SetName(*token.Properties.GetName())
		updateToken.SetProperties(*updateProp)
		msg := fmt.Sprintf("delete all scopes from Token: %s", *token.Id)
		if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
			return err
		}
		_, err = c.ContainerRegistryServices.Token().Delete(tokenId, regId)
		if err != nil {
			return err
		}
		_, _, err = c.ContainerRegistryServices.Token().Put(tokenId, *updateToken, regId)
		if err != nil {
			return err
		}

		return nil
	}

	id, err := c.Command.Command.Flags().GetInt(FlagScopeId)
	if err != nil {
		return err
	}
	id--

	updateToken := sdkgo.NewPutTokenInputWithDefaults()
	updateProp := sdkgo.NewPostTokenPropertiesWithDefaults()

	scopes := *token.Properties.GetScopes()
	scopes = append(scopes[:id], scopes[id+1:]...)
	updateProp.SetExpiryDate(*token.Properties.GetExpiryDate())
	updateProp.SetStatus(*token.Properties.GetStatus())
	updateProp.SetName(*token.Properties.GetName())
	updateProp.SetScopes(scopes)
	updateToken.SetProperties(*updateProp)

	msg := fmt.Sprintf("delete scope %d from Token: %s", id+1, *token.Id)
	if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
		return err
	}
	_, err = c.ContainerRegistryServices.Token().Delete(tokenId, regId)
	if err != nil {
		return err
	}
	_, _, err = c.ContainerRegistryServices.Token().Put(tokenId, *updateToken, regId)
	if err != nil {
		return err
	}

	return nil
}

func PreCmdTokenScopesDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{FlagRegId, FlagTokenId, FlagScopeId},
		[]string{FlagRegId, FlagTokenId, constants.ArgAll},
	)
}
