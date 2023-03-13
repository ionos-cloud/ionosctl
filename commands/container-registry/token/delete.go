package token

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "delete",
			Aliases:    []string{"d", "del", "rm"},
			ShortDesc:  "Delete a token",
			LongDesc:   "Delete a token from a registry",
			Example:    "ionosctl container-registry token delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]",
			PreCmdRun:  PreCmdDeleteToken,
			CmdRun:     CmdDeleteToken,
			InitClient: true,
		},
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all tokens from all registries")
	cmd.AddStringFlag(FlagRegId, "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddBoolFlag(FlagAllTokens, "", false, "Delete all tokens from a registry")
	cmd.AddStringFlag(FlagTokenId, "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdDeleteToken(c *core.CommandConfig) error {
	allFlag := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))
	if !allFlag {
		regId := viper.GetString(core.GetFlagName(c.NS, FlagRegId))
		allTokensFlag := viper.GetBool(core.GetFlagName(c.NS, FlagAllTokens))
		if !allTokensFlag {
			tokenId := viper.GetString(core.GetFlagName(c.NS, FlagTokenId))
			msg := fmt.Sprintf("delete Token: %s", tokenId)
			if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
				return err
			}
			_, err := c.ContainerRegistryServices.Token().Delete(tokenId, regId)
			if err != nil {
				return err
			}
			return nil
		}
		tokens, _, err := c.ContainerRegistryServices.Token().List(regId)
		if err != nil {
			return err
		}
		for _, token := range *tokens.GetItems() {
			msg := fmt.Sprintf("delete Token: %s", *token.Id)
			if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
				return err
			}
			_, err := c.ContainerRegistryServices.Token().Delete(*token.Id, regId)
			if err != nil {
				return err
			}
		}

		return nil
	}
	regs, _, err := c.ContainerRegistryServices.Registry().List("")
	if err != nil {
		return err
	}
	for _, reg := range *regs.GetItems() {
		tokens, _, err := c.ContainerRegistryServices.Token().List(*reg.Id)
		if err != nil {
			return err
		}
		for _, token := range *tokens.GetItems() {
			msg := fmt.Sprintf("delete Token: %s", *token.Id)
			if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
				return err
			}
			_, err := c.ContainerRegistryServices.Token().Delete(*token.Id, *reg.Id)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func PreCmdDeleteToken(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{FlagRegId, FlagTokenId},
		[]string{FlagRegId, FlagAllTokens},
		[]string{constants.ArgAll},
	)
}
