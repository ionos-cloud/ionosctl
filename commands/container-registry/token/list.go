package token

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "list",
			Aliases:    []string{"l"},
			ShortDesc:  "List all Registries",
			LongDesc:   "List all managed container registries for your account",
			Example:    "ionosctl container-registry registry list",
			PreCmdRun:  PreCmdListToken,
			CmdRun:     CmdListToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		"name", "n", "",
		"Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive",
	)

	cmd.AddBoolFlag("all", "a", false, "List all tokens, including expired ones")
	cmd.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
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

func CmdListToken(c *core.CommandConfig) error {
	allFlag := viper.GetBool(core.GetFlagName(c.NS, "all"))
	if !allFlag {
		id := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
		tokens, _, err := c.ContainerRegistryServices.Token().List(id)
		if err != nil {
			return err
		}
		list := tokens.GetItems()
		return c.Printer.Print(getTokenPrint(nil, c, list, false))
	}
	var list []ionoscloud.TokenResponse
	regs, _, err := c.ContainerRegistryServices.Registry().List("")
	if err != nil {
		return err
	}
	for _, reg := range *regs.GetItems() {
		tokens, _, err := c.ContainerRegistryServices.Token().List(*reg.Id)
		if err != nil {
			return err
		}
		list = append(list, *tokens.GetItems()...)
	}
	return c.Printer.Print(getTokenPrint(nil, c, &list, false))
}

func PreCmdListToken(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{"registry-id"},
		[]string{constants.ArgAll},
	)
}
