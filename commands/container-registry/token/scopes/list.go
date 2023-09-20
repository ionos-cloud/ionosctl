package scopes

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
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
			ShortDesc:  "Get a token scopes",
			LongDesc:   "Use this command to list all scopes of a token of a container registry.",
			Example:    "ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]",
			PreCmdRun:  PreCmdTokenScopesList,
			CmdRun:     CmdGetTokenScopesList,
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
		FlagTokenId, func(cobracmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(viper.GetString(core.GetFlagName(cmd.NS, FlagRegId))), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allScopeCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allScopeCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdGetTokenScopesList(c *core.CommandConfig) error {
	reg_id := viper.GetString(core.GetFlagName(c.NS, FlagRegId))
	token_id := viper.GetString(core.GetFlagName(c.NS, FlagTokenId))

	token, _, err := c.ContainerRegistryServices.Token().Get(token_id, reg_id)
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

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", allScopeJSONPaths, scopes,
		tabheaders.GetHeaders(allScopeCols, defaultScopeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func PreCmdTokenScopesList(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagTokenId, FlagRegId)
	if err != nil {
		return err
	}
	return nil
}
