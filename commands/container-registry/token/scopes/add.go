package scopes

import (
	"context"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"

	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
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

	cmd.AddStringFlag(FlagRegId, "r", "", "Registry ID", core.RequiredFlagOption())
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

	cmd.AddStringFlag(FlagName, "n", "", "Scope name", core.RequiredFlagOption())
	cmd.AddStringFlag("type", "y", "", "Scope type", core.RequiredFlagOption())
	cmd.AddStringSliceFlag("actions", "a", []string{}, "Scope actions", core.RequiredFlagOption())

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allScopeCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allScopeCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func PreCmdTokenScopesAdd(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagTokenId, FlagRegId, FlagName, "actions", "type")
	if err != nil {
		return err
	}
	return nil
}

func CmdTokenScopesAdd(c *core.CommandConfig) error {
	var scope sdkgo.Scope
	var err error

	reg_id, err := c.Command.Command.Flags().GetString(FlagRegId)
	if err != nil {
		return err
	}

	token_id, err := c.Command.Command.Flags().GetString(FlagTokenId)
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString(FlagName)
	if err != nil {
		return err
	}

	actions, err := c.Command.Command.Flags().GetStringSlice("actions")
	if err != nil {
		return err
	}

	scope_type, err := c.Command.Command.Flags().GetString("type")
	if err != nil {
		return err
	}

	scope.SetName(name)
	scope.SetActions(actions)
	scope.SetType(scope_type)

	token, _, err := c.ContainerRegistryServices.Token().Get(token_id, reg_id)
	if err != nil {
		return err
	}

	updateToken := sdkgo.NewPatchTokenInput()
	updateToken.SetExpiryDate(*token.Properties.GetExpiryDate())
	updateToken.SetStatus(*token.Properties.GetStatus())
	scopes := *token.Properties.GetScopes()
	scopes = append(scopes, scope)
	updateToken.SetScopes(scopes)

	tokenUp, _, err := c.ContainerRegistryServices.Token().Patch(token_id, *updateToken, reg_id)
	if err != nil {
		return err
	}

	return c.Printer.Print(getTokenScopePrint(nil, c, &tokenUp, true))
}
