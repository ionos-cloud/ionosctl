package token

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
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

	cmd.AddBoolFlag("all", "a", false, "Delete all tokens from all registries")
	cmd.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddBoolFlag("all-tokens", "", false, "Delete all tokens from a registry")
	cmd.AddStringFlag("token-id", "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"token-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	allFlag := viper.GetBool(core.GetFlagName(c.NS, "all"))
	if !allFlag {
		reg_id := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
		allTokensFlag := viper.GetBool(core.GetFlagName(c.NS, "all-tokens"))
		if !allTokensFlag {
			token_id := viper.GetString(core.GetFlagName(c.NS, "token-id"))
			msg := fmt.Sprintf("delete Token: %s", token_id)
			if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
				return err
			}
			_, err := c.ContainerRegistryServices.Token().Delete(token_id, reg_id)
			if err != nil {
				return err
			}
			return nil
		}
		tokens, _, err := c.ContainerRegistryServices.Token().List(reg_id)
		if err != nil {
			return err
		}
		for _, token := range *tokens.GetItems() {
			msg := fmt.Sprintf("delete Token: %s", *token.Id)
			if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
				return err
			}
			_, err := c.ContainerRegistryServices.Token().Delete(*token.Id, reg_id)
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
		[]string{"registry-id", "token-id"},
		[]string{"registry-id", "all-tokens"},
		[]string{constants.ArgAll},
	)
}
