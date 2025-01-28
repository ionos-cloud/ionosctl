package token

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(AllTokenCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return AllTokenCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdListToken(c *core.CommandConfig) error {
	allFlag := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if !allFlag {
		id := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))

		tokens, _, err := c.ContainerRegistryServices.Token().List(id)
		if err != nil {
			return err
		}
		out, err := jsontabwriter.GenerateOutput(
			"items", jsonpaths.ContainerRegistryToken, tokens, tabheaders.GetHeadersAllDefault(AllTokenCols, cols),
		)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
		return nil
	}

	var list = make([]containerregistry.TokensResponse, 0)

	regs, _, err := c.ContainerRegistryServices.Registry().List("")
	if err != nil {
		return err
	}

	if items, ok := regs.GetItemsOk(); ok && items != nil {
		for _, reg := range *items {
			tokens, _, err := c.ContainerRegistryServices.Token().List(*reg.Id)
			if err != nil {
				return err
			}

			list = append(list, tokens)
		}
	}

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.ContainerRegistryToken, list, tabheaders.GetHeadersAllDefault(AllTokenCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func PreCmdListToken(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagRegistryId},
		[]string{constants.ArgAll},
	)
}
