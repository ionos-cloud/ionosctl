package scopes

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenScopesAddCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "token",
			Resource:  "scope",
			Verb:      "add",
			Aliases:   []string{"a", "ad"},
			ShortDesc: "Add scopes to a token",
			LongDesc:  "Use this command to add scopes to a token of a container registry.",
			Example: "ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]" +
				" --actions [SCOPE-ACTIONS], --type [SCOPE-TYPE]",
			PreCmdRun:  PreCmdTokenScopesAdd,
			CmdRun:     CmdTokenScopesAdd,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID", core.RequiredFlagOption(),
	)
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Scope name", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagType, "y", "", "Scope type", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagType,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"repository", "namespace", "registry"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringSliceFlag(FlagActions, "a", []string{}, "Scope actions", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagActions,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"*", "push", "pull", "delete", "read", "write", "list"}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allScopeCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allScopeCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func PreCmdTokenScopesAdd(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(
		c.Command, c.NS, FlagTokenId, constants.FlagRegistryId, constants.FlagName, FlagActions, FlagType,
	)
	if err != nil {
		return err
	}
	return nil
}

func CmdTokenScopesAdd(c *core.CommandConfig) error {
	var scope containerregistry.Scope
	var err error

	regId, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	tokenId, err := c.Command.Command.Flags().GetString(FlagTokenId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	actions, err := c.Command.Command.Flags().GetStringSlice(FlagActions)
	if err != nil {
		return err
	}

	scopeType, err := c.Command.Command.Flags().GetString(FlagType)
	if err != nil {
		return err
	}

	scope.SetName(name)
	scope.SetActions(actions)
	scope.SetType(scopeType)

	token, _, err := c.ContainerRegistryServices.Token().Get(tokenId, regId)
	if err != nil {
		return err
	}

	updateToken := containerregistry.NewPatchTokenInput()
	if token.Properties.ExpiryDate != nil {
		updateToken.SetExpiryDate(token.Properties.GetExpiryDate())
	}

	if token.Properties.Status != nil {
		updateToken.SetStatus(token.Properties.GetStatus())
	}

	scopes := token.Properties.GetScopes()
	scopes = append(scopes, scope)

	updateToken.SetScopes(scopes)

	tokenUp, _, err := c.ContainerRegistryServices.Token().Patch(tokenId, *updateToken, regId)
	if err != nil {
		return err
	}

	scopesConverted := resource2table.ConvertContainerRegistryTokenScopesToTable(tokenUp.Properties.Scopes)

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutputPreconverted(
		token.Properties.Scopes, scopesConverted, tabheaders.GetHeadersAllDefault(allScopeCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
